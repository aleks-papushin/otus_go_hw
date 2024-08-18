package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

const (
	defaultChSize = 100
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

type UserEmail struct {
	Email string
}

type DomainStat map[string]int

var targetDomain string

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	targetDomain = domain
	scanner := bufio.NewScanner(r)
	uBytesCh := make(chan []byte, defaultChSize)

	// get users line by line
	go func() {
		for scanner.Scan() {
			uBytes := scanner.Bytes()
			uBytesCopy := make([]byte, len(uBytes))
			copy(uBytesCopy, uBytes)
			uBytesCh <- uBytesCopy
		}
		close(uBytesCh)
	}()

	// get emails
	domainsCh := make(chan string, defaultChSize)
	go func() {
		for uBytes := range uBytesCh {
			var userEmail UserEmail
			err := jsoniter.Unmarshal(uBytes, &userEmail)
			if err != nil {
				fmt.Println(err)
				return
			}
			domainsCh <- strings.ToLower(strings.Split(userEmail.Email, "@")[1])
		}
		close(domainsCh)
	}()

	// get matched mails
	matchedDomainsCh := make(chan string, defaultChSize)
	go func() {
		for e := range domainsCh {
			if strings.Split(e, ".")[1] == targetDomain {
				matchedDomainsCh <- e
			}
		}
		close(matchedDomainsCh)
	}()

	// count stats
	stats := make(DomainStat)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for sd := range matchedDomainsCh {
			num := stats[sd]
			num++
			stats[sd] = num
		}
	}()

	wg.Wait()
	return stats, nil
}
