package creditstore

import (
	"github.com/dgraph-io/badger/v3"
	"io"
	"runtime"
)

type storage struct {
	db *badger.DB
}

func (s storage) Begin(update bool) (Tx, error) {
	txn := s.db.NewTransaction(update)
	return newTx(txn), nil
}

func (s storage) Close() error {
	_ = s.Sync()
	_ = s.db.Close()
	return nil
}

func (s storage) Sync() error {
	return s.db.Sync()
}

func (s storage) DropAccount(account Account) error {
	return s.db.DropPrefix([]byte(account.String() + "."))
}

func (s storage) Backup(dst io.Writer) error {
	_, err := s.db.Backup(dst, 0)
	return err
}

func (s storage) Clean() error {
	return s.db.Flatten(runtime.NumCPU())
}
func NewStorage(path string) (CreditStore, error) {
	conf := badger.DefaultOptions(path)
	if len(path) == 0 {
		conf.InMemory = true
	}
	conf = conf.WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(conf)
	if err != nil {
		return nil, err
	}
	err = db.VerifyChecksum()
	if err != nil {
		return nil, err
	}
	s := storage{}
	s.db = db
	return s, nil
}
