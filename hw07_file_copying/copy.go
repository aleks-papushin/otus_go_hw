package main

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	wg                       = sync.WaitGroup{}
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if ok := strings.HasSuffix(fromPath, ".txt"); !ok {
		return ErrUnsupportedFile
	}

	fromFile, err := os.Open(fromPath)
	defer closeFile(fromFile)
	if err != nil {
		return err
	}
	fi, _ := fromFile.Stat()
	if offset > fi.Size() {
		return ErrOffsetExceedsFileSize
	}

	dir := filepath.Dir(toPath)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer closeFile(toFile)

	ch := make(chan int64)
	go runProgressBar(ch, offset, limit, fi.Size())
	ch <- 0

	buffer := make([]byte, 1024)
	readFrom := offset
	var read int
	var readTotal int64

	for {
		read, err = fromFile.ReadAt(buffer, readFrom)
		if err != nil && err != io.EOF {
			return err
		}
		readTotal += int64(read)

		ch <- readTotal

		if limit > 0 && readTotal > limit {
			read -= int(readTotal - limit)
		}
		_, err = toFile.WriteAt(buffer[:read], readFrom-offset)
		if err != nil {
			return err
		}

		if read != len(buffer) {
			break
		}

		readFrom += int64(read)
	}

	close(ch)
	wg.Wait()
	return nil
}

func closeFile(fromFile *os.File) {
	func(fromFile *os.File) {
		if err := fromFile.Close(); err != nil {
			fmt.Println(fmt.Errorf("error %w while closing file '%s'", err, fromFile.Name()))
		}
	}(fromFile)
}

func runProgressBar(ch <-chan int64, o, l, fs int64) {
	wg.Add(1)
	var bytesToWrite int64
	if l == 0 {
		bytesToWrite = fs - o
	} else {
		bytesToWrite = limit
	}

	bar := pb.StartNew(int(bytesToWrite))

	for readTotal := range ch {
		bar.SetCurrent(readTotal)
		time.Sleep(time.Millisecond * 400) // throttle in purpose to see progress bar animation in terminal
	}
	defer func() {
		bar.SetCurrent(bytesToWrite)
		bar.Finish()
		wg.Done()
	}()
}
