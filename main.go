package main

import (
	"fmt"
	"iht/pan"
)

func main() {
	cookie := "UID=1143614_I1_1733727252; CID=2aa893259e85a1927ac346b38cddc66d; SEID=e89346545ce9ff8767ca0757edd117ffe6fc65cefb5e654929f1749db8253e9a6133aa0027c7a4c911108a75743d80b1678751a6618d267bb3c4e9f1; KID=170f5d1d5321115a0fc8b326d2274519"
	dirid := "64456581448"

	expid, err := pan.ExportDir(cookie, dirid)
	fmt.Println(expid, err)
	edr, err := pan.ExportResult(cookie, expid)
	fmt.Println(edr, err)
	edi, err := pan.ExportPath(cookie, edr.PickCode)
	fmt.Println(edi, err)

	data, err := pan.ExportDownload(edi.Cookie, edi.FileUrl)
	fmt.Println(string(data), err)

	// data, err := pan.ExportDownload(cookie, "")
	// fmt.Println(string(data), err)
}
