package main

import (
	"fmt"

	"github.com/ebiiim/conbukun/pkg/ao/data"
)

func main() {
	for k, v := range data.Maps {
		fmt.Printf("%+v\n", k)
		fmt.Printf("\t%+v\n", v)
	}

	fmt.Print("\n\n\n")

	m, _ := data.GetMapDataFromName("Qiitun-Vietis")
	fmt.Printf("%+v\n", m)
}
