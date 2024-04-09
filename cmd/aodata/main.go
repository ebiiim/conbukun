package main

import (
	"fmt"

	"github.com/ebiiim/conbukun/pkg/ao/data"
)

func main() {
	c := 0
	for k, v := range data.Maps {
		fmt.Printf("%04d | %+v\n", c, k)
		fmt.Printf("\t%+v\n", v)
		c++
	}

	fmt.Print("\n\n")
	fmt.Printf("loaded %d maps\n", c)
	fmt.Print("\n\n")

	m, _ := data.GetMapDataFromName("Qiitun-Vietis")
	fmt.Printf("%+v\n", m)
}
