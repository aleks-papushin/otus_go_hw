package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
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
	domainCh := make(chan string)
	scanWG.Add(1)
	go func() {
		defer scanWG.Done()
		for scanner.Scan() {
			d := strings.ToLower(string(scanner.Bytes()))
			domainCh <- d
		}
	}()

	matchedCh := make(chan string)
	reWG := sync.WaitGroup{}
	reWG.Add(1)
	go func() {
		defer reWG.Done()
		for dString := range domainCh {
			dString := dString
			reWG.Add(1)
			go func() {
				defer reWG.Done()
				matched, err := regexp.Match("\\."+domain, []byte(dString))
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
			num := result[d]
			num++
			result[d] = num
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
