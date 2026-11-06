package websearch

import "io"

// Large model/browser files can live on a portable USB/exFAT disk. Bulk I/O
// avoids tens of thousands of tiny reads across that filesystem boundary.
func copyBulk(dst io.Writer, src io.Reader) (int64, error) {
	return io.CopyBuffer(struct{ io.Writer }{dst}, struct{ io.Reader }{src}, make([]byte, 4<<20))
}
