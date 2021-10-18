package creditstore

import (
	"github.com/itsabgr/go-handy"
	"math/big"
)

type credit struct {
	_   handy.NoCopy
	int *big.Int
}

func (c *credit) Neg() {
	c.int = new(big.Int).Neg(c.int)
}

func (c *credit) Encode() []byte {
	b, err := c.int.GobEncode()
	handy.Throw(err)
	return b
}

func (c *credit) Add(c2 Credit) {
	c.int = new(big.Int).Add(c.int, c2.Int())
}

func (c *credit) Sub(c2 Credit) {
	c.int = new(big.Int).Sub(c.int, c2.Int())
}

func (c *credit) String() string {
	return c.int.String()
}

func (c *credit) Sign() int {
	return c.int.Sign()
}

func (c *credit) Int() *big.Int {
	return new(big.Int).Set(c.int)
}

func (c *credit) Cmp(credit2 Credit) int {
	return credit2.Int().Cmp(c.Int())
}
func AsCredit(n *big.Int) Credit {
	c := new(credit)
	c.int = n
	return c
}
func NewCredit(n *big.Int) Credit {
	c := new(credit)
	c.int = new(big.Int).Set(n)
	return c
}

func DecodeCredit(b []byte) (Credit, error) {
	n := new(big.Int)
	err := n.GobDecode(b)
	if err != nil {
		return nil, err
	}
	return AsCredit(n), nil
}
