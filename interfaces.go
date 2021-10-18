package creditstore

import (
	"io"
	"math/big"
)

type Credit interface {
	Encode() []byte
	Add(Credit)
	Sub(Credit)
	Neg()
	String() string
	Sign() int
	Int() *big.Int
	Cmp(credit2 Credit) int
}
type Account interface {
	String() string
}
type CreditStore interface {
	Begin(update bool) (Tx, error)
	Close() error
	Sync() error
	DropAccount(account Account) error
	Backup(dst io.Writer) error
	Clean() error
}
type Lock interface {
	Encode() []byte
	Account() Account
	ID() string
	Credit() Credit
}
type Tx interface {
	Lock(account Account, amount Credit, deadline uint64) (Lock, error)
	Transfer(lock Lock, to Account, amount Credit) error
	Unlock(Lock) error
	Sum(account Account) (Credit, error)
	Closed() bool
	IsMoreThan(Account, Credit) (bool, error)
	Close() error
	Commit() error
}
