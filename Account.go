package creditstore

import (
	"errors"
	"strings"
)

type account string

func (a account) String() string {
	return string(a)
}

func NewAccount(name string) Account {
	if strings.Contains(name, ".") {
		panic(errors.New(`account name contains "."`))
	}
	a := account(name)
	return a
}
