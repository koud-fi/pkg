package feistel

import (
	"crypto/sha512"
	"encoding/binary"
)

// Encrypt applies the Feistel network on a 64-bit input using the provided round keys.
func Encrypt(plain uint64, keys []uint32) uint64 {
	var (
		left  = uint32(plain >> 32)
		right = uint32(plain & 0xFFFFFFFF)
	)
	for _, key := range keys {
		left, right = right, left^roundFunction(right, key)
	}
	return (uint64(left) << 32) | uint64(right)
}

// Decrypt reverses the Feistel network operation by applying the round keys in reverse.
func Decrypt(cipher uint64, keys []uint32) uint64 {
	var (
		left  = uint32(cipher >> 32)
		right = uint32(cipher & 0xFFFFFFFF)
	)
	for i := len(keys) - 1; i >= 0; i-- {
		key := keys[i]
		left, right = right^roundFunction(left, key), left
	}
	return (uint64(left) << 32) | uint64(right)
}

// roundFunction is a simple non-linear function used in each round.
// It rotates the input r left by 5 bits and then adds the key.
func roundFunction(r, key uint32) uint32 {
	rotated := (r << 5) | (r >> (32 - 5))
	return rotated + key
}

// KeysFromString generates a slice of round keys from a abitrary string.
// The input string is decoded and used to seed a sequence of keys.
func KeysFromString(s string) (keys [16]uint32) {
	hash := sha512.Sum512([]byte(s))
	for i := range keys {
		keys[i] = binary.BigEndian.Uint32(hash[i*4 : (i+1)*4])
	}
	return keys
}
