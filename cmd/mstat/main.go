package main

import "github.com/bhavya-dang/mstat/internal/cli"

var version = "dev"

func main() {
	cli.SetVersion(version)
	cli.Execute()
}
