package util

import (
	"math/rand"
	"strings"
	"time"
)

// 创建一个独立的随机数生成器
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// 字母表
const alphabet = "abcdefghijklmnopqrstuvwxyz"

// 生成min,max间的一个随机整数
func RandomInt(min, max int64) int64 {
	return min + seededRand.Int63n(max-min+1)
}

// 生成一个随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[seededRand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[seededRand.Intn(n)]
}
