package flags

import (
	"my-ls-1/internal/models"
	"os"
	"strings"
)

func Parse() (models.Flags, []string) {
	var f models.Flags
	var paths []string

	args := os.Args[1:]

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			for _, ch := range arg[1:] {
				switch ch {
				case 'l':
					f.Long = true
				case 'a':
					f.All = true
				case 'r':
					f.Reverse = true
				case 'R':
					f.Recursive = true
				}
			}
		} else {
			paths = append(paths, arg)
		}
	}
	return f, paths
}
