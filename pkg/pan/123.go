package pan

import (
	"encoding/json"
	"fmt"
	"iht/utils"
	"time"

	"github.com/xshrim/gol/tk"
)

var header123 = map[string]string{
	"Platform":     "open_platform",
	"Content-Type": "application/json",
}

type P123 struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
	Token  string `json:"token"`
	Expiry string `json:"expiry"`
}

type FileObj struct {
	FileId       int64  `json:"fileId"`
	Filename     string `json:"filename"`
	Type         int    `json:"type"`
	Size         int64  `json:"size"`
	Etag         string `json:"etag"`
	Status       int    `json:"status"`
	ParentFileId int    `json:"parentFileId"`
	Category     int    `json:"category"`
}

func (c *P123) GetToken() (string, string, error) {
	if c.Token != "" && c.Expiry != "" {
		expiry, err := time.Parse(time.RFC3339, c.Expiry)
		if err == nil && time.Now().Before(expiry) {
			return c.Token, c.Expiry, nil
		}
	}

	urlstr := "https://open-api.123pan.com/api/v1/access_token"

	header := make(map[string]string)
	for key, value := range header123 {
		header[key] = value
	}

	payload := map[string]string{
		"clientID":     c.Id,
		"clientSecret": c.Secret,
	}

	jsBytes, _ := json.Marshal(payload)

	data, _, err := utils.Post(urlstr, string(jsBytes), time.Second*30, header)
	if err != nil {
		return "", "", err
	}

	if val := tk.Jsquery(string(data), ".message"); val != nil {
		if v, ok := val.(string); ok && v == "ok" {
			if token, ok := tk.Jsquery(string(data), ".data.accessToken").(string); ok {
				c.Token = token
			}

			if expiry, ok := tk.Jsquery(string(data), ".data.expiredAt").(string); ok {
				c.Expiry = expiry
			}
		}
	} else {
		return "", "", fmt.Errorf("get token failed: %v", string(data))
	}

	return c.Token, c.Expiry, nil
}

func (c *P123) FetchList(cid string) ([]FileObj, error) {
	if cid == "" {
		cid = "0"
	}
	var flist []FileObj

	last := 0
	for last != -1 {
		urlstr := fmt.Sprintf("https://open-api.123pan.com/api/v2/file/list?parentFileId=%s&limit=100", cid)
		if last != 0 {
			urlstr += fmt.Sprintf("&lastFileId=%d", last)
			time.Sleep(1 * time.Second)
		}

		header := make(map[string]string)
		for key, value := range header123 {
			header[key] = value
		}

		token, _, err := c.GetToken()
		if err != nil {
			return nil, err
		}
		header["Authorization"] = "Bearer " + token

		data, _, err := utils.Get(urlstr, time.Second*30, header)
		if err != nil {
			return nil, err
		}

		var tlist []FileObj
		if val := tk.Jsquery(string(data), ".message"); val != nil {
			if v, ok := val.(string); ok && v == "ok" {
				last = int(tk.Jsquery(string(data), ".data.lastFileId").(float64))
				items := tk.Jsquery(string(data), ".data.fileList")
				if err := json.Unmarshal([]byte(tk.Jsonify(items)), &tlist); err != nil {
					return nil, fmt.Errorf("list dir failed: %v", err)
				}
				flist = append(flist, tlist...)
			} else {
				return nil, fmt.Errorf("list dir failed: %v", tk.Jsquery(string(data), ".error"))
			}
		} else {
			return nil, fmt.Errorf("list dir failed: %v", string(data))
		}
	}

	return flist, nil

}
