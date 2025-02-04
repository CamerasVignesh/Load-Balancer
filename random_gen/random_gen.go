package randomgen

import (
	"math/rand/v2"
	"strconv"
)

func RandRange() string {
    return strconv.Itoa(rand.IntN(10000-8081) + 8081)
}
