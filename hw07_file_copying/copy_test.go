package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	defer func() {
		if err := os.RemoveAll("tmp"); err != nil {
			fmt.Println(fmt.Errorf("error attempting to delete path '%s': %w", "tmp", err))
		}
	}()
	fromPath := "testdata/input.txt"
	tests := []struct {
		name          string
		fromPath      string
		outPath       string
		compareToPath string
		offset        int64
		limit         int64
	}{
		{"offset 0 limit 0", fromPath, "tmp/out0-0.txt", "testdata/out_offset0_limit0.txt", 0, 0},
		{"offset 0 limit 10", fromPath, "tmp/out0-10.txt", "testdata/out_offset0_limit10.txt", 0, 10},
		{"offset 0 limit 1000", fromPath, "tmp/out0-1000.txt", "testdata/out_offset0_limit1000.txt", 0, 1000},
		{"offset 0 limit 10000", fromPath, "tmp/out0-10000.txt", "testdata/out_offset0_limit10000.txt", 0, 10000},
		{"offset 100 limit 1000", fromPath, "tmp/out100-1000.txt", "testdata/out_offset100_limit1000.txt", 100, 1000},
		{"offset 6000 limit 1000", fromPath, "tmp/out6000-1000.txt", "testdata/out_offset6000_limit1000.txt", 6000, 1000},
		{"offset 1000 limit 0", fromPath, "tmp/out1000-0.txt", "testdata/out_offset1000_limit0.txt", 1000, 0},
	}

	for _, tc := range tests {
		fmt.Printf("running test case %s\n", tc.name)

		err := Copy(tc.fromPath, tc.outPath, tc.offset, tc.limit)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		require.True(t, sameFiles(tc.outPath, tc.compareToPath), "files are not the same")

		if err = os.Remove(tc.outPath); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func sameFiles(path1, path2 string) bool {
	f1Hash := getFileHash(path1)
	f2Hash := getFileHash(path2)

	hash1Str := hex.EncodeToString(f1Hash.Sum(nil))
	hash2Str := hex.EncodeToString(f2Hash.Sum(nil))

	return hash1Str == hash2Str
}

func getFileHash(path string) hash.Hash {
	f, _ := os.Open(path)
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}(f)

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Println(err)
	}

	return h
}
