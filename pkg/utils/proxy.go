package utils

import (
	"math/rand"
	"time"
	"unsafe"
)

func RandomDNS(dns []string) string {
	return dns[rand.Intn(len(dns))]
}

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
