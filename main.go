package main

import (
	"fmt"
	"myM3u8/parse"
	"os"
)

func main() {
	f, _ := os.Open("./index.list")

	//	fts,_:=io.ReadAll(f)
	a, e := parse.Parse(f)
	if e != nil {
		fmt.Println(e)
		os.Exit(101)
	}
	fmt.Println(a)
}
