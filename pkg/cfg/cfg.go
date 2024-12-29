package cfg

import (
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Cookie string `json:"cookie"`
	Cid    string `json:"cid"`
	Cpath  string `json:"cpath"`
	Base   string `json:"base"`
	Tree   string `json:"tree"`
	Url    string `json:"url"`
	Prefix string `json:"prefix"`
}

var Conf Config

func init() {
	if dataBytes, err := os.ReadFile("config.yaml"); err == nil {
		_ = yaml.Unmarshal(dataBytes, &Conf)
	}

	cookie := flag.String("cookie", "", "115 cookie")
	cid := flag.String("cid", "", "115 directory/file id")
	cpath := flag.String("cpath", "", "115 directory/file path")
	base := flag.String("base", "media", "media home directory")
	tree := flag.String("tree", "", "directory tree path")
	prefix := flag.String("prefix", "", "media url prefix path")
	url := flag.String("url", "", "media url")
	flag.Parse()

	if *cookie != "" {
		Conf.Cookie = *cookie
	}
	if *cid != "" {
		Conf.Cid = *cid
	}
	if *cpath != "" {
		Conf.Cpath = *cpath
	}
	if *base != "" {
		Conf.Base = *base
	}
	if *tree != "" {
		Conf.Tree = *tree
	}
	if *url != "" {
		Conf.Url = *url
	}
	if *prefix != "" {
		Conf.Prefix = *prefix
	}

	if Conf.Cookie == "" || (Conf.Cid == "" && Conf.Cpath == "") || Conf.Base == "" || Conf.Url == "" || Conf.Prefix == "" {
		panic("config variables <cookie | cid | cpath | base | url | prefix> are not set")
	}
}
