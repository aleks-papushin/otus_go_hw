package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"sync"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(scanBetweenSubstrings(byte('@'), byte('"')))
	scanWG := sync.WaitGroup{}
	domainCh := make(chan []byte)
	scanWG.Add(1)
	go func() {
		defer scanWG.Done()
		for scanner.Scan() {
			dBytes := scanner.Bytes()
			dBytesCopy := make([]byte, len(dBytes))
			copy(dBytesCopy, dBytes)
			scanWG.Add(1)
			go func() {
				defer scanWG.Done()
				domainCh <- bytesToLower(dBytesCopy)
			}()
		}
	}()

	matchedCh := make(chan []byte)
	reWG := sync.WaitGroup{}
	reWG.Add(1)
	go func() {
		defer reWG.Done()
		for dString := range domainCh {
			dString := dString
			reWG.Add(1)
			go func() {
				defer reWG.Done()
				matched, err := regexp.Match("\\."+domain, dString)
				if !matched || err != nil {
					return
				}
				matchedCh <- dString
			}()
		}
	}()

	result := make(DomainStat)
	resWG := sync.WaitGroup{}
	resWG.Add(1)
	go func() {
		defer resWG.Done()
		for d := range matchedCh {
			dString := string(d)
			num := result[dString]
			num++
			result[dString] = num
		}
	}()
	scanWG.Wait()
	close(domainCh)
	reWG.Wait()
	close(matchedCh)
	resWG.Wait()
	return result, nil
}

func scanBetweenSubstrings(start, end byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		startIndex := bytes.IndexByte(data, start)
		if startIndex == -1 {
			if atEOF {
				return 0, nil, nil
			}
			return 0, nil, nil
		}
		startIndex++
		endIndex := bytes.IndexByte(data[startIndex:], end)
		if endIndex == -1 {
			if atEOF {
				return 0, nil, nil
			}
			return 0, nil, nil
		}
		endIndex += startIndex
		return endIndex + 1, data[startIndex:endIndex], nil
	}
}

func bytesToLower(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return b
}
