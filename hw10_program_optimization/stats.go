//go:generate easyjson -all stats.go

package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

// easyjson:json
type UserEmail struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	ds := make(DomainStat)
	var errors []error

	for scanner.Scan() {
		uBytes := scanner.Bytes()
		ubCopy := make([]byte, len(uBytes))
		copy(ubCopy, uBytes)

		var uEmail UserEmail
		err := easyjson.Unmarshal(ubCopy, &uEmail)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		email := strings.ToLower(uEmail.Email)

		if strings.HasSuffix(email, "."+domain) {
			if strings.Contains(email, "@") {
				userDomain := strings.Split(email, "@")[1]
				num := ds[userDomain]
				num++
				ds[userDomain] = num
			}
		}
	}

	if len(errors) > 0 {
		return ds, fmt.Errorf("unmarshal errors: %v", errors)
	}

	return ds, nil
}
