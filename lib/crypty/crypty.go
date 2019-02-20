// Package Crypty is a tiny cryptography lib
// It was made just to avoid the "/action?somethind=2" replacing with some wierd numbers
package Crypty

import "math/rand"

// Generate a randon Int Hash
var hash = rand.Intn(672389128)

// Get the value to crypto
func Crypto(v int) int {
	return v * hash
}

// Decrypt the value
func Decrypto(v int) int {
	return v / hash
}
