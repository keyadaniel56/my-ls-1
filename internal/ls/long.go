package ls

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"my-ls-1/internal/models"
)

func getStat(info os.FileInfo) (*syscall.Stat_t, bool) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	return stat, ok
}

// formatMode converts os.FileMode to the exact ls -l format string.
func formatMode(m os.FileMode) string {
	buf := [10]byte{}

	// file type character
	switch {
	case m&os.ModeSymlink != 0:
		buf[0] = 'l'
	case m&os.ModeDir != 0:
		buf[0] = 'd'
	case m&os.ModeNamedPipe != 0:
		buf[0] = 'p'
	case m&os.ModeSocket != 0:
		buf[0] = 's'
	case m&os.ModeDevice != 0:
		if m&os.ModeCharDevice != 0 {
			buf[0] = 'c'
		} else {
			buf[0] = 'b'
		}
	default:
		buf[0] = '-'
	}

	// permission bits
	const rwx = "rwxrwxrwx"
	perm := m.Perm()
	for i := 0; i < 9; i++ {
		if perm&(1<<uint(8-i)) != 0 {
			buf[1+i] = rwx[i]
		} else {
			buf[1+i] = '-'
		}
	}

	// setuid / setgid / sticky
	if m&os.ModeSetuid != 0 {
		if buf[3] == 'x' {
			buf[3] = 's'
		} else {
			buf[3] = 'S'
		}
	}
	if m&os.ModeSetgid != 0 {
		if buf[6] == 'x' {
			buf[6] = 's'
		} else {
			buf[6] = 'S'
		}
	}
	if m&os.ModeSticky != 0 {
		if buf[9] == 'x' {
			buf[9] = 't'
		} else {
			buf[9] = 'T'
		}
	}

	return string(buf[:])
}

func PrintLong(entries []models.FileEntry) {
	type row struct {
		mode  string
		nlink string
		uname string
		gname string
		size  string
		mtime string
		name  string
	}

	maxNlink := 1
	maxSize := 1
	maxUser := 1
	maxGroup := 1

	rows := make([]row, 0, len(entries))

	for _, e := range entries {
		info := e.Info
		stat, ok := getStat(info)
		if !ok {
			continue
		}

		uname := strconv.Itoa(int(stat.Uid))
		gname := strconv.Itoa(int(stat.Gid))
		if u, err := user.LookupId(uname); err == nil {
			uname = u.Username
		}
		if g, err := user.LookupGroupId(gname); err == nil {
			gname = g.Name
		}

		nlink := strconv.FormatUint(uint64(stat.Nlink), 10)
		mtime := formatTime(info.ModTime())
		name := e.Name

		// symlink arrow — info here is already from Lstat (see readEntries)
		if info.Mode()&os.ModeSymlink != 0 {
			if target, err := os.Readlink(e.Path); err == nil {
				name = e.Name + " -> " + target
			}
		}

		// size: for block/char devices show "major, minor"
		var size string
		mode := info.Mode()
		if mode&os.ModeDevice != 0 {
			major := uint64(stat.Rdev>>8) & 0xfff
			minor := uint64(stat.Rdev) & 0xff
			size = fmt.Sprintf("%3d, %5d", major, minor)
		} else {
			size = strconv.FormatInt(info.Size(), 10)
		}

		r := row{
			mode:  formatMode(mode),
			nlink: nlink,
			uname: uname,
			gname: gname,
			size:  size,
			mtime: mtime,
			name:  name,
		}
		rows = append(rows, r)

		if len(nlink) > maxNlink {
			maxNlink = len(nlink)
		}
		if len(size) > maxSize {
			maxSize = len(size)
		}
		if len(uname) > maxUser {
			maxUser = len(uname)
		}
		if len(gname) > maxGroup {
			maxGroup = len(gname)
		}
	}

	for _, r := range rows {
		fmt.Printf("%s %*s %-*s %-*s %*s %s %s\n",
			r.mode,
			maxNlink+1, r.nlink,
			maxUser, r.uname,
			maxGroup, r.gname,
			maxSize, r.size,
			r.mtime,
			colorName(r.name),
		)
	}
}

func formatTime(t time.Time) string {
	now := time.Now()
	if now.Sub(t) > 6*30*24*time.Hour || t.After(now) {
		return t.Format("Jan _2  2006")
	}
	return t.Format("Jan _2 15:04")
}
