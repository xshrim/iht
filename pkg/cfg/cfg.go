package cfg

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	P115 P115           `yaml:"p115"`
	P123 P123           `yaml:"p123"`
	Flow map[string]any `yaml:"flow"`
}

type P115 struct {
	Cookie  string `json:"cookie"`
	Cid     string `json:"cid"`
	Cpath   string `json:"cpath"`
	Fpath   string `json:"fpath"`
	Library string `json:"library"`
	Url     string `json:"url"`
	Prefix  string `json:"prefix"`
}

type P123 struct {
	Id     string `json:"id"`
	Secret string `json:"secret"`
	Cid    string `json:"cid"`
	Cpath  string `json:"cpath"`
}

var Conf Config

func init() {
	if dataBytes, err := os.ReadFile("config.yaml"); err == nil {
		_ = yaml.Unmarshal(dataBytes, &Conf)
	}

	// cookie := flag.String("cookie", "", "115 cookie")
	// cid := flag.String("cid", "", "115 directory/file id")
	// cpath := flag.String("cpath", "", "115 directory/file path")
	// base := flag.String("base", "media", "media home directory")
	// tree := flag.String("tree", "", "directory tree path")
	// prefix := flag.String("prefix", "", "media url prefix path")
	// url := flag.String("url", "", "media url")
	// flag.Parse()

	// if *cookie != "" {
	// 	Conf.P115.Cookie = *cookie
	// }
	// if *cid != "" {
	// 	Conf.P115.Cid = *cid
	// }
	// if *cpath != "" {
	// 	Conf.P115.Cpath = *cpath
	// }
	// if *base != "" {
	// 	Conf.P115.Base = *base
	// }
	// if *tree != "" {
	// 	Conf.P115.Tree = *tree
	// }
	// if *url != "" {
	// 	Conf.P115.Url = *url
	// }
	// if *prefix != "" {
	// 	Conf.P115.Prefix = *prefix
	// }

	// if Conf.P115.Cookie == "" || (Conf.P115.Cid == "" && Conf.P115.Cpath == "") || Conf.P115.Base == "" || Conf.P115.Url == "" || Conf.P115.Prefix == "" {
	// 	panic("config variables <cookie | cid | cpath | base | url | prefix> for pan 115 are not set")
	// }
}
