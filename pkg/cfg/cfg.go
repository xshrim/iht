package cfg

import (
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Cookie string `json:"cookie"`
	Did    string `json:"did"`
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
	did := flag.String("did", "", "115 directory id")
	base := flag.String("base", "media", "media home directory")
	tree := flag.String("tree", "", "directory tree path")
	prefix := flag.String("prefix", "", "media url prefix path")
	url := flag.String("url", "", "media url")
	flag.Parse()

	if *cookie != "" {
		Conf.Cookie = *cookie
	}
	if *did != "" {
		Conf.Did = *did
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

	if Conf.Cookie == "" || Conf.Did == "" || Conf.Base == "" || Conf.Url == "" || Conf.Prefix == "" {
		// print("config variables <cookie | did | base | url> are not set")
		panic("config variables <cookie | did | base | url | prefix> are not set")
	}
}
