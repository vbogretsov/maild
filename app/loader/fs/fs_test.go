package fs_test

import (
	"testing"

	"github.com/vbogretsov/maild/app"
	"github.com/vbogretsov/maild/app/loader/fs"
)

func TestArgumentErrorIfNoSuchFile(t *testing.T) {
	ld, _ := fs.New("/tmp")
	_, err := ld.Load("en", "example")
	if err != nil {
		switch err.(type) {
		case app.ArgumentError:
		default:
			t.Error("unexpected error type")
		}
	} else {
		t.Error("expected error but got <nil>")
	}
}
