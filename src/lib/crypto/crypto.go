package crypto

import (
	"crypto/rand"
	"crypto/sha512"

	"golang.org/x/crypto/pbkdf2"
)

func StrengthenMasterHash(masterHash []byte, salt []byte) []byte {
	return pbkdf2.Key(masterHash, salt, 150000, 512/8, sha512.New)
}

func RandomBytes(byteLength int) []byte {
	bytes := make([]byte, byteLength)
	rand.Read(bytes)
	return bytes
}
