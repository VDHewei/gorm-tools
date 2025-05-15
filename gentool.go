package main

import "github.com/VDHewei/gorm-tools/pkg/core"

func main() {
	cli := core.New()
	if cli != nil && !cli.PrintHelp() {
		cli.Execute()
	}
}
