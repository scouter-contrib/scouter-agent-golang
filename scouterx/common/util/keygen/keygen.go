package keygen

import (
	"math/rand"
	"sync"
	"time"
)

var once sync.Once

// Keygen is a KeyGen struct type variable
var keygen *KeyGen

// A KeyGen is a source of random numbers
type KeyGen struct {
	sync.Mutex
	random *rand.Rand
}

// GetInstance returns single new Keygen
func getInstance() *KeyGen {
	once.Do(
		func() {
			keygen = new(KeyGen)
			keygen.random = rand.New(rand.NewSource(time.Now().Unix()))

		})
	return keygen
}

// Next returns a non-negative pseudo-random 63-bit integer
func Next() int64 {
	keygen.Lock()
	key := keygen.random.Int63()
	keygen.Unlock()
	return key
}

func init() {
	getInstance()
}
