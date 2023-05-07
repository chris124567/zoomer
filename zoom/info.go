package zoom

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func httpGet(client *http.Client, url string, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header = headers
	return client.Do(request)
}

func (session *ZoomSession) generateSignature(meetingNumber string) string {
	ts := time.Now().Unix()

	header := []byte(`{"alg":"HS256","typ":"JWT"}`)
	payload := []byte(fmt.Sprintf(`{"sdkKey":"%s","iat":%d,"exp":%d,"mn":%s,"role":0}`, session.ZoomJwtApiKey, ts, ts+1800, meetingNumber))
	message := base64.URLEncoding.EncodeToString(header) + "." + base64.URLEncoding.EncodeToString(payload)

	h := hmac.New(sha256.New, []byte(session.ZoomJwtApiSecret))
	h.Write([]byte(message))
	return message + "." + base64.URLEncoding.EncodeToString(h.Sum(nil))
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
	values.Set("cv", "2.12.0")
	values.Set("proxy", "1")
	values.Set("sdkOrigin", "aHR0cDovL2xvY2FsaG9zdDo5OTk5")
	values.Set("tk", "")
	values.Set("ztk", "")
	values.Set("sdkUrl", "aHR0cDovL2xvY2FsaG9zdDo5OTk5L21lZXRpbmcuaHRtbA")
	values.Set("captcha", "")
	values.Set("captchaName", "")
	values.Set("suid", "")
	values.Set("callback", "axiosJsonpCallback1")
	values.Set("signatureType", "sdk")

	response, err := httpGet(session.httpClient, fmt.Sprintf("https://zoom.us/api/v1/wc/info?%s", values.Encode()), httpHeaders())
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	getStringInBetweenTwoString := func(str []byte, startS []byte, endS []byte) []byte {
		s := bytes.Index(str, startS)
		if s == -1 {
			return nil
		}
		newS := str[s+len(startS):]

		e := bytes.Index(newS, endS)
		if e == -1 {
			return nil
		}
		return newS[:e]
	}

	if err = json.Unmarshal(getStringInBetweenTwoString(data, []byte("osJsonpCallback1("), []byte(")")), &meetingInfo); err != nil {
		return nil, "", err
	}

	var cookieString string
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

func (session *ZoomSession) getRwgPingData(meetingInfo *MeetingInfo, pingRwcServer *RwgInfo) (*RwgInfo, error) {
	headers := httpHeaders()
	headers.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := httpGet(session.httpClient, fmt.Sprintf("https://%s/wc/ping/%s?ts=%d&auth=%s&rwcToken=%s&dmz=1", pingRwcServer.Rwg, meetingInfo.Result.MeetingNumber, meetingInfo.Result.Ts, meetingInfo.Result.Auth, pingRwcServer.RwcAuth), headers)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var rwgPingInfo RwgInfo
	if err := json.NewDecoder(response.Body).Decode(&rwgPingInfo); err != nil {
		return nil, err
	}
	return &rwgPingInfo, nil
}
