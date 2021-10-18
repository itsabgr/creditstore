package creditstore

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
	"math/big"
	"runtime"
	"time"
)

type tx struct {
	tx     *badger.Txn
	closed bool
}

func (t tx) Closed() bool {
	return t.closed
}

func (t tx) Close() error {
	t.tx.Discard()
	t.closed = true
	return nil
}

func (t tx) Commit() error {
	t.closed = true
	return t.tx.Commit()
}

func (t tx) set(account, index string, c Credit, deadline uint64) error {
	entry := badger.NewEntry([]byte(account+"."+index), c.Encode())
	entry.ExpiresAt = deadline
	return t.tx.SetEntry(entry)
}
func (t tx) exists(account, index string) (bool, error) {
	item, err := t.tx.Get([]byte(account + "." + index))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return false, nil
		}
		return false, err
	}
	if item.IsDeletedOrExpired() {
		return false, nil
	}
	return true, nil
}
func (t tx) get(account, index string) (Credit, error) {
	item, err := t.tx.Get([]byte(account + "." + index))
	if err != nil {
		return nil, err
	}
	if item.IsDeletedOrExpired() {
		return nil, badger.ErrKeyNotFound
	}
	val, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}
	return DecodeCredit(val)
}
func (t tx) delete(account, index string) error {
	return t.tx.Delete([]byte(account + "." + index))
}

func (t tx) inc(account string, inc Credit) error {
	c, err := t.get(account, "")
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
		c = NewCredit(big.NewInt(0))
	}
	c.Add(inc)
	if c.Sign() == 0 {
		return t.delete(account, "")
	}
	return t.set(account, "", c, 0)
}

func (t tx) Lock(account Account, amount Credit, deadline uint64) (Lock, error) {
	if amount.Sign() != 1 {
		return nil, errors.New("non positive lock")
	}
	if deadline < uint64(time.Now().Unix()) {
		return nil, errors.New("past deadline")
	}
	var lock Lock
	for {
		lock = randLock()
		exists, err := t.exists(account.String(), lock.String())
		if err != nil {
			return nil, err
		}
		if exists {
			runtime.Gosched()
			continue
		}
		break
	}
	amount.Neg()
	return lock, t.set(account.String(), lock.String(), amount, deadline)
}
func (t tx) Transfer(from Account, lock Lock, to Account, amount Credit) error {
	max, err := t.get(from.String(), lock.String())
	if err != nil {
		return err
	}
	if max.Cmp(amount) == -1 {
		return errors.New("overflow")
	}
	err = t.inc(to.String(), amount)
	if err != nil {
		return err
	}
	amount.Neg()

	err = t.inc(from.String(), amount)
	if err != nil {
		_ = t.Close()
		return err
	}

	err = t.Unlock(from, lock)
	if err != nil {
		_ = t.Close()
		return err
	}
	return nil
}

func (t tx) Unlock(a Account, l Lock) error {
	return t.delete(a.String(), l.String())
}

func (t tx) Sum(account Account) (Credit, error) {
	prefix := []byte(account.String() + ".")
	iter := t.tx.NewIterator(badger.IteratorOptions{
		AllVersions: false,
		Prefix:      prefix,
	})
	defer iter.Close()
	total := NewCredit(big.NewInt(0))
	for iter.Rewind(); iter.ValidForPrefix(prefix); iter.Next() {
		item := iter.Item()
		if item.IsDeletedOrExpired() {
			continue
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return total, err
		}
		c, err := DecodeCredit(val)
		if err != nil {
			return total, err
		}
		total.Add(c)
	}
	return total, nil
}

func (t tx) IsMoreThan(account Account, max Credit) (bool, error) {
	prefix := []byte(account.String() + ".")
	iter := t.tx.NewIterator(badger.IteratorOptions{
		AllVersions: false,
		Prefix:      prefix,
	})
	defer iter.Close()
	total := NewCredit(big.NewInt(0))
	for iter.Rewind(); iter.ValidForPrefix(prefix); iter.Next() {
		item := iter.Item()
		if item.IsDeletedOrExpired() {
			continue
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return false, err
		}
		c, err := DecodeCredit(val)
		if err != nil {
			return false, err
		}
		total.Add(c)
		if total.Cmp(max) >= 0 {
			return true, nil
		}
	}
	return false, nil
}

func newTx(txn *badger.Txn) Tx {
	return &tx{txn, false}
}
