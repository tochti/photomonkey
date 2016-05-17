package bot

import (
	"io/ioutil"
	"os"
	"testing"
)

func Test_sha1OfFile(t *testing.T) {
	fh, err := ioutil.TempFile(".", "tmp")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fh.Name())

	_, err = fh.Write([]byte("testfile blob"))
	if err != nil {
		t.Fatal(err)
	}

	hash, err := sha1OfFile(fh.Name())
	if err != nil {
		t.Fatal(err)
	}

	expect := "a4cc75d3e7e32c015e192e7cd9b44334498ce1f9"
	if hash != expect {
		t.Fatalf("Expect %v was %v", expect, hash)
	}
}
