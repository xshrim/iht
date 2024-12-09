package pan

import (
	"fmt"
	"iht/utils"
	"net/url"
	"time"
)

var headers = map[string]string{

	"Origin":             "https://webapi.115.com",
	"Referer":            "https://webapi.115.com/",
	"Sec-Fetch-Dest":     "empty",
	"Sec-Fetch-Mode":     "cors",
	"Sec-Fetch-Site":     "same-site",
	"sec-ch-ua":          `"Not/A)Brand";v="99", "Google Chrome";v="115", "Chromium";v="115"`,
	"sec-ch-ua-mobile":   "?0",
	"sec-ch-ua-platform": "Linux",
	"User-Agent":         "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
}

func ExportDir(cookie, dirid string) (string, error) {

	urlstr := "https://webapi.115.com/files/export_dir"

	formValues := url.Values{}
	formValues.Set("file_ids", dirid)
	formValues.Set("target", "U_1_0")
	formValues.Set("layer_limit", "25")
	formDataStr := formValues.Encode()

	header := make(map[string]string)
	for key, value := range headers {
		header[key] = value
	}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	header["Cookie"] = cookie

	data, err := utils.Post(urlstr, formDataStr, time.Second*30, header)
	if err != nil {
		return "", err
	}

	fmt.Println(string(data), err)
	return "", nil
}
