package testing

import (
	"math/rand"
	"sync"
	"time"
)

var NATIVE_ACTOR_USER_URL = "native_actor_user:80"

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
var seedMutex sync.Mutex

func RandomString(length int) string {
	seedMutex.Lock()
	defer seedMutex.Unlock()

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
