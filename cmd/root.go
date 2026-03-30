package cmd

import (
	"fmt"
	"my-ls-1/internal/ls"
	"my-ls-1/internal/models"
)

func Run(flags models.Flags, paths []string) {
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, path := range paths {
		if len(paths) > 1 {
			fmt.Printf("%s:\n", path)
		}

		err := ls.List(path, flags)
		if err != nil {
			fmt.Println("error:", err)
		}
	}
}