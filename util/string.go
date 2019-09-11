package util

import (
	"math/rand"
	"strconv"
	"time"
)

//RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//NonceStr 用时间戳生成随机字符串
func NonceStr() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}
