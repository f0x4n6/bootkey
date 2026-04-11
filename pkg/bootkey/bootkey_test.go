package bootkey

import (
	"encoding/hex"
	"path/filepath"
	"testing"
)

const bk = "13d20976d63ea5e836036ec8bc68d6eb"

var path = filepath.Join("testdata", "SYSTEM")

func TestReadFile(t *testing.T) {
	t.Run("Test ReadFile", func(t *testing.T) {
		key, err := ReadFile(path)

		if err != nil {
			t.Fatalf("ReadFile: %v", err)
		}

		if hex.EncodeToString(key) != bk {
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
