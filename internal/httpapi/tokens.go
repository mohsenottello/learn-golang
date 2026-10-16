package httpapi

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// newToken returns a random token to send to the user and the hash to store.
// Only the hash is kept in the database, so a DB leak doesn't expose live links.
func newToken() (raw, hash string) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	raw = hex.EncodeToString(b)
	return raw, hashToken(raw)
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
