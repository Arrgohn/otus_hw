package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	user := User{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if strings.Contains(user.Email, "."+domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
