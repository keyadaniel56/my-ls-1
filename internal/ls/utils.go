package ls

import "my-ls-1/internal/models"


func PrintSimple(entries []models.FileEntry){
	for _,e:=range entries{
		print(e.Name +" ")
	}
	println()
}