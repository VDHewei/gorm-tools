package main

import "github.com/VDHewei/gorm-tools/pkg/core"

func main() {
	cli := core.New()
	if !cli.PrintHelp() && !cli.PrintVersion() {
		cli.Execute()
	}
}
