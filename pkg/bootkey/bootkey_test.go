package bootkey

import (
	"path/filepath"
	"testing"
)

const bk = "\x13\xD2\x09\x76\xD6\x3E\xA5\xE8\x36\x03\x6E\xC8\xBC\x68\xD6\xEB"

var path = filepath.Join("testdata", "SYSTEM")

func TestReadFile(t *testing.T) {
	t.Run("Test ReadFile", func(t *testing.T) {
		key, err := ReadFile(path)

		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}

		if string(key) != bk {
			t.Fatal("bootkey invalid")
		}
	})
}

func BenchmarkReadFile(b *testing.B) {
	b.Run("Benchmark ReadFile", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if _, err := ReadFile(path); err != nil {
				b.Fatalf("ReadFile: %v", err)
			}
		}
	})
}
