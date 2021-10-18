package creditstore

import (
	"errors"
	"strings"
)

type account string

func (a account) String() string {
	return string(a)
}

func NewAccount(name string) (Account, error) {
	if strings.Contains(name, ".") {
		return nil, errors.New(`account name contains "."`)
	}
	a := account(name)
	return a, nil
}
