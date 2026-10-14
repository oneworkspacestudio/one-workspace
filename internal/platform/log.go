package platform

import (
	"os"
	"path/filepath"
	"sync"
)

type Log struct {
	mu   sync.Mutex
	Path string
}

func (l *Log) Write(b []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e := os.MkdirAll(filepath.Dir(l.Path), 0700); e != nil {
		return 0, e
	}
	if info, e := os.Stat(l.Path); e == nil && info.Size()+int64(len(b)) > 5<<20 {
		_ = os.Remove(l.Path + ".1")
		_ = os.Rename(l.Path, l.Path+".1")
	}
	f, e := os.OpenFile(l.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if e != nil {
		return 0, e
	}
	defer f.Close()
	return f.Write(b)
}
