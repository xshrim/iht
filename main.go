package main

import (
	"fmt"
	"iht/pkg/cron"
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
	// pan.Go()

	crontab()
}
