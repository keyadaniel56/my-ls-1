package models

import "os"

type FileEntry struct {
	Name  string
	Path  string
	Info  os.FileInfo
	IsDir bool
}
