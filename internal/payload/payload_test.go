package payload

import "testing"

func TestEmbeddedPayloadCatalog(t *testing.T) {
	f, e := files.Open("bundle/payload.zip")
	if e != nil {
		t.Skip("source-only checks; no bundled binaries")
	}
	_ = f.Close()
	z, c, e := Open()
	if e != nil {
		t.Fatal(e)
	}
	if len(c.Apps) != 4 || len(z.File) < 5 {
		t.Fatal("incomplete offline catalog")
	}
}
