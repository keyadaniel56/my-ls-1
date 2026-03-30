package main

import (
	"my-ls-1/cmd"
	"my-ls-1/internal/flags"
)

func main() {
	f, paths := flags.Parse()
	cmd.Run(f, paths)
}