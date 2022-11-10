package lib

import (
	"crypto/rand"
	"crypto/sha512"

	"golang.org/x/crypto/pbkdf2"
)

func StrengthenMasterHash(masterHash []byte) ([]byte, []byte) {
	randomSalt := RandomBytes(16)
	return pbkdf2.Key(masterHash, randomSalt, 150000, 512/8, sha512.New), randomSalt
}

func RandomBytes(byteLength int) []byte {
	bytes := make([]byte, byteLength)
	rand.Read(bytes)
	return bytes
}
