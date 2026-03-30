package ls

import (
	"fmt"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"my-ls-1/internal/models"
)

func PrintLong(entries []models.FileEntry) {
	for _, e := range entries {
		info := e.Info

		stat := info.Sys().(*syscall.Stat_t)

		u, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
		g, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

		fmt.Printf("%s %d %s %s %d %s %s\n",
			info.Mode().String(),
			stat.Nlink,
			u.Username,
			g.Name,
			info.Size(),
			formatTime(info.ModTime()),
			e.Name,
		)
	}
}

func formatTime(t time.Time) string {
	return t.Format("Jan _2 15:04")
}