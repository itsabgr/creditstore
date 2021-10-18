package creditstore

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/itsabgr/go-handy"
	"io"
	"time"
)

type lock struct {
	id     string
	from   Account
	credit Credit
}

func (l lock) ID() string {
	return l.id
}

func (l lock) Account() Account {
	return l.from
}

func (l lock) Credit() Credit {
	return l.credit
}

type lockSchema struct {
	Account string `json:"account"`
	Credit  string `json:"credit"`
	ID      string `json:"id"`
}

func (l lock) Encode() []byte {
	b, err := json.Marshal(lockSchema{
		Account: l.Account().String(),
		Credit:  l.Credit().String(),
		ID:      l.id,
	})
	handy.Throw(err)
	return b
}

func DecodeLock(b []byte) (Lock, error) {
	schema := new(lockSchema)
	err := json.Unmarshal(b, schema)
	if err != nil {
		return nil, err
	}
	lock := lock{id: schema.ID}
	lock.from, err = NewAccount(schema.Account)
	if err != nil {
		return nil, err
	}

	lock.credit, err = NewCreditFromString(schema.Credit)
	if err != nil {
		return nil, err
	}
	return lock, nil
}

func randLock(a Account, c Credit) Lock {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[:8], uint64(time.Now().UTC().UnixNano()))
	_, err := io.ReadFull(rand.Reader, b[8:])
	handy.Throw(err)
	return lock{hex.EncodeToString(b), a, c}
}
