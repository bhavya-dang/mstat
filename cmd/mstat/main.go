package main

import (
	"fmt"
	"os"

	"github.com/bhavya-dang/mstat/internal/version"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println(version.Version)
		os.Exit(0)
	}
	fmt.Println("mstat — work in progress")
}
