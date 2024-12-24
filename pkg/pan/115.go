package pan

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"iht/pkg/cfg"
	"iht/utils"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xshrim/gol"
	"github.com/xshrim/gol/tk"
	"golang.org/x/text/encoding/unicode"
)

var mediaExt = []string{
	".mp4",
	".mkv",
	".avi",
	".mov",
	".flv",
	".wmv",
	".mp3",
	".wav",
	".flac",
	".aac",
	".m4a",
	".ogg",
	".webm",
	".ts",
	".m3u8",
	".flv",
	".f4v",
	".rmvb",
	".rm",
	".3gp",
	".wmv",
	".asf",
	".m4v",
	".m4p",
	".m4b",
	".m4r",
	".m4v",
	".m4a",
}

var captionExt = []string{
	".srt",
	".ass",
	".ssa",
	".vtt",
	".sbv",
	".stl",
	".dfxp",
	".ttml",
	".sub",
	".idx",
}

var ignoreExt = []string{
	".jpg",
	".jpeg",
	".png",
	".gif",
	".bmp",
	".webp",
	".ico",
	".svg",
	".tif",
	".tiff",
	".psd",
	".heic",
	".heif",
	".dng",
	".cr2",
	".nef",
	".orf",
	".arw",
	".txt",
	".nfo",
	".url",
	".htm",
	".html",
}

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

func contains(elems []string, s string) bool {
	s = strings.TrimLeft(s, "./")
	for _, elem := range elems {
		if elem == s || strings.HasPrefix(filepath.Dir(elem)+"/", s+"/") {
			return true
		}
	}
	return false
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

func ExportTree(cookie, dirid string) ([]byte, error) {
	expid, err := ExportDir(cookie, dirid)
	if err != nil {
		return nil, err
	}

	edr, err := ExportResult(cookie, expid)
	if err != nil {
		return nil, err
	}

	edi, err := ExportPath(cookie, edr.PickCode)
	if err != nil {
		return nil, err
	}

	data, err := ExportDownload(edi.Cookie, edi.FileUrl)
	if err != err {
		return nil, err
	}

	if err := ExportDelete(cookie, edi.FileId); err != nil {
		return nil, err
	}

	// UTF16转UTF8
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	idata, err := decoder.Bytes(data)
	if err != nil {
		return nil, err
	}

	return idata, nil
}

func ExportTreeToFile(cookie, dirid, path string) (string, error) {

	data, err := ExportTree(cookie, dirid)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, data, 0666); err != nil {
		return "", err
	}

	return filepath.Abs(path)
}

func Tree2List(data []byte) ([]string, error) {
	var list []string
	var stack []string
	layer := 0

	var reader *bufio.Reader
	file, err := os.Open(string(data))

	if err != nil {
		reader = bufio.NewReader(bytes.NewReader(data))
	} else {
		defer file.Close()
		reader = bufio.NewReader(file)
	}

	for {
		bline, _, err := reader.ReadLine()
		if err != nil {
			// if err == io.EOF {
			// 	break
			// }
			break
		}

		line := string(bline)

		if strings.HasSuffix(line, "根目录") {
			continue
		} else {
			iline := strings.Split(line, "|-")
			head := iline[0]
			name := iline[1]
			curlayer := strings.Count(head, "|")
			if curlayer <= layer {
				list = append(list, strings.Join(stack, "/"))
				stack = stack[:layer-1]
				if curlayer < layer {
					list = append(list, strings.Join(stack, "/"))
					stack = stack[:curlayer-1]
				}
			}
			stack = append(stack, name)
			layer = curlayer
		}
	}

	return list, nil
}

func Tree2Lib(iurl, prefix, dir string, data []byte) error {
	list, err := Tree2List(data)
	if err != nil {
		return err
	}

	var items []string

	for _, item := range list {
		ext := filepath.Ext(item)
		if utils.Contains(mediaExt, ext) && !utils.Contains(ignoreExt, ext) {
			fpath := filepath.Join(dir, fmt.Sprintf("%s.strm", item))
			items = append(items, fpath)
			if _, err := os.Stat(fpath); os.IsNotExist(err) {
				// create directory if necessary
				if err := os.MkdirAll(filepath.Dir(fpath), 0777); err != nil {
					return err
				}
				urlstr, _ := url.JoinPath(iurl, "/d", prefix, item)
				// writefile with override
				if err := os.WriteFile(fpath, []byte(urlstr), 0666); err != nil {
					return err
				}
			}
		}
	}

	return CleanNotExist(dir, items)
}

func CleanNotExist(dir string, list []string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// only remove directories and strm file that not exist in list
		if !contains(list, path) && (info.IsDir() || filepath.Ext(path) == ".strm") {
			ftype := "file"
			if info.IsDir() {
				ftype = "directory"
			}

			if err := os.RemoveAll(path); err != nil {
				gol.Warnf("Redundant %s %s remove failed: %v", ftype, path, err)
			} else {
				gol.Infof("Redundant %s %s remove succeed", ftype, path)
			}
		}

		return nil
	})
}

func Go() {
	gol.Info("Start to export directory tree")
	// pan.ExportTreeToFile(cfg.Conf.Cookie, cfg.Conf.Did, cfg.Conf.Tree)
	data, err := ExportTree(cfg.Conf.Cookie, cfg.Conf.Did)
	if err != nil {
		gol.Errorf("Export directory tree failed: %v", err)
		return
	}
	gol.Info("Export directory tree succeed")

	gol.Info("Start to convert directory tree to library")
	if err := Tree2Lib(cfg.Conf.Url, cfg.Conf.Prefix, cfg.Conf.Base, data); err != nil {
		gol.Errorf("Convert directory tree to library failed: %v", err)
		return
	}
	gol.Info("Convert directory tree to library succeed")
}
