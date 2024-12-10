package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/xshrim/gol"
)

func bashEscape(str string) string {
	return `'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}

func Timing(t time.Time) int64 {
	return time.Since(t).Milliseconds()
}

// Http2Curl returns a CurlCommand corresponding to an http.Request
func Http2Curl(req *http.Request) (string, error) {
	var command []string
	if req.URL == nil {
		return "", fmt.Errorf("getCurlCommand: invalid request, req.URL is nil")
	}

	command = append(command, "curl")

	schema := req.URL.Scheme
	requestURL := req.URL.String()
	if schema == "" {
		schema = "http"
		if req.TLS != nil {
			schema = "https"
		}
		requestURL = schema + "://" + req.Host + req.URL.Path
	}

	if schema == "https" {
		command = append(command, "-k")
	}

	command = append(command, "-X", bashEscape(req.Method))

	if req.Body != nil {
		var buff bytes.Buffer
		_, err := buff.ReadFrom(req.Body)
		if err != nil {
			return "", fmt.Errorf("getCurlCommand: buffer read from body error: %w", err)
		}
		// reset body for potential re-reads
		req.Body = ioutil.NopCloser(bytes.NewBuffer(buff.Bytes()))
		if len(buff.String()) > 0 {
			bodyEscaped := bashEscape(buff.String())
			command = append(command, "-d", bodyEscaped)
		}
	}

	var keys []string

	for k := range req.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		command = append(command, "-H", bashEscape(fmt.Sprintf("%s: %s", k, strings.Join(req.Header[k], " "))))
	}

	command = append(command, bashEscape(requestURL))

	command = append(command, "--compressed")

	return strings.Join(command, " "), nil
}

func GenCookie(method, url string, fd map[string]string) string {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range fd {
		w.WriteField(k, v)
	}
	w.Close()
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	cookies := ""
	for _, cookie := range resp.Cookies() {
		cookies += cookie.Name + "=" + cookie.Value + ";"
	}
	cookies = strings.TrimSuffix(cookies, ";")
	return cookies
}

func Request(method, url string, headers map[string]any, payload []byte) (int, []byte, error) {
	var resp *http.Response
	var err error
	// TODO 支持队列
	client := &http.Client{
		Timeout: time.Duration(30) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		if strings.Contains(url, ":443") {
			url = "https://" + url
		} else {
			url = "http://" + url
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return 0, nil, err
	}

	if method == "POST" {
		if _, ok := headers["Content-Type"]; !ok {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	for k, v := range headers {
		vstr, _ := v.(string)
		req.Header.Add(k, vstr)
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(payload))
	resp, err = client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	code := resp.StatusCode

	data, err := ioutil.ReadAll(resp.Body)
	return code, data, err
}

func TimingRequest(url, header, payload string, timeout int) (latency int64, msg string, sc int, err error) {
	defer func(t time.Time) {
		latency = Timing(t)
	}(time.Now())

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	method := "GET"
	if payload != "" {
		method = "POST"
	}

	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return
	}

	for _, kvstr := range strings.Split(header, "|") {
		kv := strings.Split(kvstr, ": ")
		if len(kv) == 2 {
			req.Header.Add(kv[0], kv[1])
		}
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	msgBytes, err := ioutil.ReadAll(resp.Body)
	sc = resp.StatusCode
	msg = string(msgBytes)

	return
}

func Get(url string, timeout time.Duration, header map[string]string) (data []byte, cookie string, err error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	gol.Debug(Http2Curl(req))

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// cookie = ""
	// for _, ck := range resp.Cookies() {
	// 	cookie += fmt.Sprintf("%s=%s;", ck.Name, ck.Value)
	// }
	// cookie = strings.TrimSuffix(cookie, ";")
	for _, ck := range resp.Header.Values("Set-Cookie") {
		cookie = ck
	}

	data, err = ioutil.ReadAll(resp.Body)

	return
}

func Post(url, payload string, timeout time.Duration, header map[string]string) (data []byte, cookie string, err error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	gol.Debug(Http2Curl(req))

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	for _, ck := range resp.Header.Values("Set-Cookie") {
		cookie = ck
	}

	data, err = ioutil.ReadAll(resp.Body)

	return
}

func Header(url, payload string, timeout time.Duration) (header http.Header, err error) {

	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	header = resp.Header

	return
}
