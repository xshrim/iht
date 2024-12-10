package pan

import (
	"encoding/json"
	"fmt"
	"iht/utils"
	"net/url"
	"time"

	"github.com/xshrim/gol/tk"
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

type ExportDirResult struct {
	ExportId string `json:"export_id"`
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	PickCode string `json:"pick_code"`
}

type ExportDownloadInfo struct {
	UserId   int64  `json:"user_id"`
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileSize string `json:"file_size"`
	FileUrl  string `json:"file_url"`
	PickCode string `json:"pick_code"`
	Cookie   string `json:"cookie"`
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

	data, _, err := utils.Post(urlstr, formDataStr, time.Second*30, header)
	if err != nil {
		return "", err
	}

	if val := tk.Jsquery(string(data), ".state"); val != nil {
		if v, ok := val.(bool); ok && v {
			if expid, ok := tk.Jsquery(string(data), ".data.export_id").(float64); ok {
				return fmt.Sprintf("%d", int64(expid)), nil
			} else {
				return "", fmt.Errorf("export dir failed: invalid export id")
			}
		} else {
			return "", fmt.Errorf("export dir failed: %v", tk.Jsquery(string(data), ".error"))
		}
	} else {
		return "", fmt.Errorf("export dir failed: %v", string(data))
	}
}

func ExportResult(cookie, expid string) (ExportDirResult, error) {
	urlstr := fmt.Sprintf("https://webapi.115.com/files/export_dir?export_id=%s", expid)

	header := make(map[string]string)
	for key, value := range headers {
		header[key] = value
	}
	header["Cookie"] = cookie

	retry := 5
	for {
		data, _, err := utils.Get(urlstr, time.Second*30, header)
		if err != nil {
			return ExportDirResult{}, err
		}

		if val := tk.Jsquery(string(data), ".state"); val != nil {
			if v, ok := val.(bool); ok && v {
				edr := ExportDirResult{}
				expdata, _ := tk.Jsquery(string(data), ".data").(any)
				if err := json.Unmarshal([]byte(tk.Jsonify(expdata)), &edr); err != nil {
					if retry <= 0 {
						return ExportDirResult{}, fmt.Errorf("get export dir result failed: %v", err)
					} else {
						retry--
						time.Sleep(time.Second * 6)
						continue
					}
				}
				return edr, nil
			} else {
				return ExportDirResult{}, fmt.Errorf("get export dir result failed: %v", tk.Jsquery(string(data), ".error"))
			}
		} else {
			return ExportDirResult{}, fmt.Errorf("get export dir result failed: %v", string(data))
		}
	}
}

func ExportPath(cookie, pickcode string) (ExportDownloadInfo, error) {
	urlstr := fmt.Sprintf("https://webapi.115.com/files/download?pickcode=%s", pickcode)

	header := make(map[string]string)
	for key, value := range headers {
		header[key] = value
	}
	header["Cookie"] = cookie

	data, respcookie, err := utils.Get(urlstr, time.Second*30, header)
	if err != nil {
		return ExportDownloadInfo{}, err
	}

	if val := tk.Jsquery(string(data), ".state"); val != nil {
		edi := ExportDownloadInfo{}
		if err := json.Unmarshal(data, &edi); err != nil {
			return ExportDownloadInfo{}, fmt.Errorf("get export download path failed: %v", err)
		}
		edi.Cookie = respcookie
		return edi, nil
	} else {
		return ExportDownloadInfo{}, fmt.Errorf("get export download path failed: %v", string(data))
	}
}

func ExportDownload(cookie, fileurl string) ([]byte, error) {
	urlstr := fileurl
	header := make(map[string]string)
	for key, value := range headers {
		header[key] = value
	}
	header["Cookie"] = cookie
	// header["Cookie"] = "44b16881664cb98f87584819b332961e=55c3aee57415e0e2031c367e652794da; expires=Tue, 10-Dec-2024 03:12:17 GMT; Max-Age=978; path=/6757ace041b94a2535948cfdc00ae76de00da427/; domain=115.com"
	// urlstr = "https://cdnfhnfile.115.com/6757ace041b94a2535948cfdc00ae76de00da427/%E6%A0%B9%E7%9B%AE%E5%BD%9520241210105216_%E7%9B%AE%E5%BD%95%E6%A0%91.txt?t=1733800337&u=1143614&s=104857600&d=974333514-b5p5a81tht66b8l0u-0&c=0&f=3&k=3c19ba0f59cd62cc145271d5c632bf60&us=1048576000&uc=10&v=1"

	data, _, err := utils.Get(urlstr, time.Second*30, header)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ExportDelete(cookie, fid string) error {
	urlstr := "https://webapi.115.com/rb/delete"

	header := make(map[string]string)
	for key, value := range headers {
		header[key] = value
	}

	header["Content-Type"] = "application/x-www-form-urlencoded"
	header["Cookie"] = cookie

	formValues := url.Values{}
	formValues.Set("pid", "0")
	formValues.Set("fid[0]", fid)
	formValues.Set("ignore_warn", "1")
	formDataStr := formValues.Encode()

	data, _, err := utils.Post(urlstr, formDataStr, time.Second*30, header)
	if err != nil {
		return fmt.Errorf("delete export dir result file failed: %v", err)
	}

	if val := tk.Jsquery(string(data), ".state"); val != nil {
		if v, ok := val.(bool); ok && v {
			return nil
		} else {
			return fmt.Errorf("delete export dir result file failed: %v", tk.Jsquery(string(data), ".error"))
		}
	} else {
		return fmt.Errorf("delete export dir result file failed: %v", string(data))
	}
}
