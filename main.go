package main

import (
	"fmt"
	"iht/cmd"
	"iht/pkg/cron"
	"iht/pkg/flow"
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

func runflow() {
	flow, err := flow.Load("./flow.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(flow.Run([]string{"月球陨落.Moonfall.2022.2160p.WEB-DL.x265.10bit.HDR.DDP5.1.Atmos-NOGRP.mkv"}))
	//fmt.Println(flow.Run("月球.mkv"))
	//fmt.Println(utils.RuneIndex("月球陨落.Moonfall.2022", 4))
}

func main() {
	// TODO crontab
	// pan.ToFile()
	// pan.Export()

	// fitem, _ := pan.Locate("/我的视频/电影/日韩/王国")
	// fmt.Println(fitem)

	// flist, _ := pan.List(fitem.Cid)
	// fmt.Println(flist)

	// attr, _ := pan.Attr(fitem.Cid)
	// fmt.Println(attr)

	// crontab()

	//runflow()

	cmd.Execute()
}
