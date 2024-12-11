package main

import (
	"bufio"
	"fmt"
	"iht/pan"
	"os"
	"strings"
)

func export() {
	cookie := "UID=1143614_I1_1733727252; CID=2aa893259e85a1927ac346b38cddc66d; SEID=e89346545ce9ff8767ca0757edd117ffe6fc65cefb5e654929f1749db8253e9a6133aa0027c7a4c911108a75743d80b1678751a6618d267bb3c4e9f1; KID=170f5d1d5321115a0fc8b326d2274519"
	dirid := "64456581448"
	fpath := "tree.txt"

	// expid, err := pan.ExportDir(cookie, dirid)
	// fmt.Println(expid, err)
	// edr, err := pan.ExportResult(cookie, expid)
	// fmt.Println(edr, err)
	// edi, err := pan.ExportPath(cookie, edr.PickCode)
	// fmt.Println(edi, err)

	// data, err := pan.ExportDownload(edi.Cookie, edi.FileUrl)
	// fmt.Println(string(data), err)

	// err = pan.ExportDelete(cookie, edi.FileId)
	// fmt.Println(err)

	fmt.Println(pan.ExportTree(cookie, dirid, fpath))

	// fmt.Println(edi.FileId)

	// data, err := pan.ExportDownload(cookie, "")
	// fmt.Println(string(data), err)
}

func parse(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var stack []string
	layer := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "|——") {
			continue
		} else {
			iline := strings.Split(line, "|-")
			head := iline[0]
			name := iline[1]
			curlayer := strings.Count(head, "|")
			if curlayer <= layer {
				fmt.Println(strings.Join(stack, "/"))
				stack = stack[:layer]
				if curlayer < layer {
					fmt.Println(strings.Join(stack, "/"))
					stack = stack[:curlayer]
				}
			}
			stack = append(stack, name)
			layer = curlayer
		}
	}
	return nil

}

func main() {
	parse("tree.txt")
}
