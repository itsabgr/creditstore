package creditstore

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"github.com/itsabgr/go-handy"
	"io"
	"time"
)

type lock []byte

func (l lock) String() string {
	return hex.EncodeToString(l)
}

func DecodeLock(str string) (Lock, error) {
	b, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return asLock(b), nil
}
func asLock(b []byte) Lock {
	return lock(b)
}

func randLock() Lock {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[:8], uint64(time.Now().UTC().UnixNano()))
	_, err := io.ReadFull(rand.Reader, b[8:])
	handy.Throw(err)
	return asLock(b)
}
