package creditstore

import (
	"fmt"
	"github.com/itsabgr/go-handy"
	"math/big"
	"testing"
	"time"
)

func TestBadger(t *testing.T) {
	cs, err := NewStorage("")
	handy.Throw(err)
	defer cs.Close()
	tx, err := cs.Begin(true)
	handy.Throw(err)
	defer tx.Close()
	ali, err := NewAccount("1")
	handy.Throw(err)
	hasan, err := NewAccount("12")
	handy.Throw(err)
	lock, err := tx.Lock(ali, NewCredit(big.NewInt(100)), uint64(time.Now().Add(10*time.Second).Unix()))
	handy.Throw(err)
	err = tx.Transfer(lock, hasan, NewCredit(big.NewInt(9)))
	handy.Throw(err)
	fmt.Println(tx.Sum(ali))
	fmt.Println(tx.Sum(hasan))
}
