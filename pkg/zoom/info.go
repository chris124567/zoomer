package zoom

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/chris124567/zoomer/internal/util"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (session *ZoomSession) generateSignature(meetingNumber string) string {
	timeStamp := (time.Now().UTC().UnixNano() / 1e6) - 30000
	timeStampStr := strconv.Itoa(int(timeStamp))

	msg := session.ZoomJwtApiKey + meetingNumber + timeStampStr + ZOOM_ROLE
	sEnc := base64.StdEncoding.EncodeToString([]byte(msg))
	h := hmac.New(sha256.New, []byte(session.ZoomJwtApiSecret))
	h.Write([]byte(sEnc))

	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	sigNatureStr := session.ZoomJwtApiKey + "." + meetingNumber + "." + timeStampStr + "." + ZOOM_ROLE + "." + hash
	return base64.StdEncoding.EncodeToString([]byte(sigNatureStr))
}

func (session *ZoomSession) GetMeetingInfoData() (*MeetingInfo, string, error) {
	var meetingInfo MeetingInfo

	// generate info url
	values := url.Values{}
	values.Set("meetingNumber", session.MeetingNumber)
	values.Set("userName", session.Username)
	values.Set("passWord", session.MeetingPassword)
	values.Set("signature", session.generateSignature(session.MeetingNumber))
	// values.Set("apiKey", ZOOM_JWT_API_KEY)
	values.Set("apiKey", session.ZoomJwtApiKey)
	values.Set("lang", "en-US")
	values.Set("userEmail", "")
	values.Set("cv", "1.8.6")
	// values.Set("cv", "1.8.5")
	values.Set("proxy", "1")
	values.Set("sdkOrigin", "aHR0cDovL2xvY2FsaG9zdDo5OTk5")
	values.Set("tk", "")
	values.Set("ztk", "")
	values.Set("sdkUrl", "aHR0cDovL2xvY2FsaG9zdDo5OTk5L21lZXRpbmcuaHRtbA")
	values.Set("captcha", "")
	values.Set("captchaName", "")
	values.Set("suid", "")
	values.Set("callback", "axiosJsonpCallback1")
	infoUrl := (&url.URL{
		Scheme:   "https",
		Host:     "zoom.us",
		Path:     "/api/v1/wc/info",
		RawQuery: values.Encode(),
	}).String()

	response, err := httpGet(session.httpClient, infoUrl, INITIAL_HEADERS)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	jsonData := util.GetStringInBetweenTwoString(string(data), "osJsonpCallback1(", ")")
	err = json.Unmarshal([]byte(jsonData), &meetingInfo)
	if err != nil {
		return nil, "", err
	}

	cookieString := ""
	for _, cookieValue := range response.Cookies() {
		// skip "cred" and empty cookies
		if cookieValue.Name == "cred" || cookieValue.Value != "" {
			cookieString += fmt.Sprintf("%s=%s; ", cookieValue.Name, cookieValue.Value)
		}
	}
	// remove trailing ";"" if it exists
	cookieString = strings.TrimSuffix(cookieString, "; ")

	return &meetingInfo, cookieString, nil
}

func getRwgPingServer(meetingInfo *MeetingInfo) *RwgInfo {
	var rwgPingInfo RwgInfo

	for key, value := range meetingInfo.Result.EncryptedRWC {
		rwgPingInfo.Rwg = key
		rwgPingInfo.RwcAuth = value
		break
	}

	return &rwgPingInfo
}

func (session *ZoomSession) getRwgPingData(meetingInfo *MeetingInfo, pingRwcServer *RwgInfo) (*RwgInfo, error) {
	var rwgPingInfo RwgInfo

	headers := util.CopyMap(INITIAL_HEADERS)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	pingUrl := (&url.URL{
		Scheme: "https",
		Host:   pingRwcServer.Rwg,
		Path:   fmt.Sprintf("/wc/ping/%s", meetingInfo.Result.MeetingNumber),
		// THE ORDER OF PARAMTERS IS VERY IMPORTANT! IT DOES NOT WORK OTHERWISE
		RawQuery: fmt.Sprintf("ts=%d&auth=%s&rwcToken=%s&dmz=1", meetingInfo.Result.Ts, meetingInfo.Result.Auth, pingRwcServer.RwcAuth),
		// RawQuery: values.Encode(),
	}).String()

	response, err := httpGet(session.httpClient, pingUrl, INITIAL_HEADERS)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&rwgPingInfo)
	if err != nil {
		return nil, err
	}

	return &rwgPingInfo, nil
}
