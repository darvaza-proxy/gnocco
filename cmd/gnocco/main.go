// Package main is the main executable
package main

import (
	"fmt"

	"darvaza.org/gnocco"
)

func main() {
	if err := gnocco.ListenAndServe(); err != nil {
		fmt.Println(err)
		return
	}
}
