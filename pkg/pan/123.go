package pan

import (
	"encoding/json"
	"fmt"
	"iht/utils"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/xshrim/gol"
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
	Path         string `json:"path"`
	Category     int    `json:"category"`
}

type FileInfo struct {
	FileID       int64  `json:"fileID"`
	Filename     string `json:"filename"`
	Type         int    `json:"type"`
	Size         int64  `json:"size"`
	Etag         string `json:"etag"`
	Status       int    `json:"status"`
	ParentFileID int    `json:"parentFileID"`
	Path         string `json:"path"`
	CreateAt     string `json:"createAt"`
	Trashed      int    `json:"trashed"`
}

func (c *P123) request(urlstr string, payload string) ([]byte, error) {
	header := make(map[string]string)
	for key, value := range header123 {
		header[key] = value
	}

	var data []byte
	var err error

	for {
		token, _, err := c.GetToken()
		if err != nil {
			return nil, err
		}
		header["Authorization"] = "Bearer " + token

		if payload == "" {
			data, _, err = utils.Get(urlstr, time.Second*30, header)
		} else {
			data, _, err = utils.Post(urlstr, payload, time.Second*30, header)
		}

		if err != nil {
			return nil, err
		}

		if fmt.Sprintf("%v", tk.Jsquery(string(data), ".code")) == "401" {
			gol.Warn("token expired, refreshing...\n")
			c.Token = ""
		} else {
			break
		}
	}

	return data, err
}

func (c *P123) GetToken() (string, string, error) {
	if c.Token != "" {
		if c.Expiry == "" {
			return c.Token, c.Expiry, nil
		}
		expiry, err := time.Parse(time.RFC3339, c.Expiry)
		if err != nil {
			return "", "", err
		}
		if time.Now().Before(expiry) {
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

func (c *P123) FetchObj(fpath string) (FileObj, error) {
	cid := "0"
	flist, err := c.FetchList(cid)
	if err != nil {
		return FileObj{}, err
	}

	fpath = strings.TrimSuffix(fpath, "/")
	if fpath == "" {
		return FileObj{FileId: 0, Filename: "根目录", Type: 1}, nil
	}
	for idx, fname := range strings.Split(fpath, "/") {
		if fname == "" {
			continue
		}

		for _, fobj := range flist {
			if fobj.Filename == fname {
				if idx == len(strings.Split(fpath, "/"))-1 {
					return fobj, nil
				} else {
					var err error
					cid = fmt.Sprintf("%d", fobj.FileId)
					flist, err = c.FetchList(cid)
					if err != nil {
						return FileObj{}, err
					}
					break
				}
			}
		}
	}

	last := strings.Split(fpath, "/")[len(strings.Split(fpath, "/"))-1]
	for _, fobj := range flist {
		if strings.Contains(fobj.Filename, last) {
			return fobj, nil
		}
	}

	return FileObj{}, fmt.Errorf("file %s not found", fpath)
}

func (c *P123) FetchPath(cid string) (string, error) {
	if cid == "0" || cid == "" {
		return "/", nil
	}

	paths := []string{}
	for {
		fstat, err := c.FetchInfo(cid)
		if err != nil {
			return "", err
		}

		paths = append([]string{fstat.Filename}, paths...)

		if fstat.ParentFileID == 0 {
			break
		}

		cid = fmt.Sprintf("%d", fstat.ParentFileID)
	}

	return "/" + strings.Join(paths, "/"), nil
}

func (c *P123) FetchInfo(cid string) (FileInfo, error) {
	if cid == "0" || cid == "" {
		return FileInfo{FileID: 0, Filename: "根目录", Type: 1, Size: 0, Etag: "", Status: 1, ParentFileID: 0, Path: "/", CreateAt: "", Trashed: 0}, nil
	}

	urlstr := fmt.Sprintf("https://open-api.123pan.com/api/v1/file/detail?fileID=%s", cid)

	data, err := c.request(urlstr, "")
	if err != nil {
		return FileInfo{}, err
	}

	var fstat FileInfo
	if val := tk.Jsquery(string(data), ".message"); val != nil {
		if v, ok := val.(string); ok && v == "ok" {
			items := tk.Jsquery(string(data), ".data")
			if err := json.Unmarshal([]byte(tk.Jsonify(items)), &fstat); err != nil {
				return FileInfo{}, fmt.Errorf("get file attribute failed: %v", err)
			}
		} else {
			return FileInfo{}, fmt.Errorf("get file attribute failed: %v", tk.Jsquery(string(data), ".message"))
		}
	} else {
		return FileInfo{}, fmt.Errorf("get file attribute failed: %v", string(data))
	}

	if err := json.Unmarshal(data, &fstat); err != nil {
		return FileInfo{}, fmt.Errorf("get file attribute failed: %v", err)
	}

	return fstat, nil
}

func (c *P123) FetchAttr(cid string) (FileInfo, error) {
	fstat, err := c.FetchInfo(cid)
	if err != nil {
		return FileInfo{}, nil
	}

	fstat.Path, _ = c.FetchPath(cid)

	return fstat, nil
}

func (c *P123) FetchList(key ...string) ([]FileObj, error) {
	// 单个参数表示获取指定目录下所有文件，两个参数表示获取指定目录下满足搜索条件的文件
	var param string
	var cpath string
	if len(key) == 0 {
		param = fmt.Sprintf("parentFileId=%s", "0")
		cpath = "/"
	} else if len(key) == 1 || key[1] == "" {
		param = fmt.Sprintf("parentFileId=%s", key[0])
		cpath, _ = c.FetchPath(key[0])
	} else {
		smode := "0"
		if len(key) > 2 {
			smode = key[2]
		}
		param = fmt.Sprintf("parentFileId=%s&searchData=%s&searchMode=%s", key[0], key[1], smode)
	}

	var flist []FileObj

	last := 0
	for last != -1 {
		urlstr := fmt.Sprintf("https://open-api.123pan.com/api/v2/file/list?%s&limit=100", param)
		if last != 0 {
			urlstr += fmt.Sprintf("&lastFileId=%d", last)
		}

		data, err := c.request(urlstr, "")
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
				return nil, fmt.Errorf("list dir failed: %v", tk.Jsquery(string(data), ".message"))
			}
		} else {
			return nil, fmt.Errorf("list dir failed: %v", string(data))
		}
	}

	var wg sync.WaitGroup
	for i := range flist {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if cpath != "" {
				flist[i].Path = filepath.Join(cpath, flist[i].Filename)
			} else {
				flist[i].Path, _ = c.FetchPath(fmt.Sprintf("%d", flist[i].FileId))
			}
		}(i)
	}
	wg.Wait()

	return flist, nil

}

func (c *P123) RenameList(rnlist []string) map[string][]string {
	urlstr := fmt.Sprintf("https://open-api.123pan.com/api/v1/file/rename")

	result := make(map[string][]string)
	result["succeed"] = []string{}
	result["failed"] = []string{}

	for _, rn := range utils.ChunkSlice(rnlist, 30) {
		payload := map[string][]string{
			"renameList": rn,
		}

		jsBytes, _ := json.Marshal(payload)

		data, err := c.request(urlstr, string(jsBytes))
		if err != nil {
			result["failed"] = append(result["failed"], rn...)
			gol.Warnf("rename file failed: %v\n", err)
		}

		if val := tk.Jsquery(string(data), ".message"); val != nil {
			if v, ok := val.(string); ok && v == "ok" {
				result["succeed"] = append(result["succeed"], rn...)
			} else {
				result["failed"] = append(result["failed"], rn...)
				gol.Warnf("rename file failed: %v\n", tk.Jsquery(string(data), ".message"))
			}
		} else {
			result["failed"] = append(result["failed"], rn...)
			gol.Warnf("rename file failed: %v\n", string(data))
		}
	}

	return result
}
