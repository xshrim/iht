package main

import (
	"fmt"
	"iht/pkg/cron"
	"iht/pkg/pan"
	"sync"
)

func crontab() {

	ctab := cron.New()

	var wg sync.WaitGroup
	wg.Add(10)

	if err := ctab.AddJob("* * * * *", func() { wg.Done(); fmt.Println("Hello, World!") }); err != nil {
		fmt.Println(err)
	}

	if err := ctab.AddJob("* * * * *", func(s string) { wg.Done() }, "param"); err != nil {
		fmt.Println(err)
	}

	ctab.RunAll()
	wg.Wait()
}

func main() {
	// TODO crontab
	// pan.ToFile()
	pan.Export()

	// fitem, _ := pan.Locate("/我的视频/电影/日韩/肮脏")
	// fmt.Println(fitem)

	// flist, _ := pan.List(fitem.Cid)
	// fmt.Println(flist)

	// attr, _ := pan.Attr(fitem.Cid)
	// fmt.Println(attr)

	// crontab()
}
