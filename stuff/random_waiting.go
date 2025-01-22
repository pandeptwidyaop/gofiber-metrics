package stuff

import (
	"math/rand"
	"time"
)

func Wait() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	time.Sleep(time.Duration(rng.Intn(1000)) * time.Millisecond)
}
