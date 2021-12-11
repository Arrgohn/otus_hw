package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	user := User{}

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		sp := strings.Split(user.Email, ".")
		dom := sp[len(sp)-1:][0]

		if dom == domain {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
