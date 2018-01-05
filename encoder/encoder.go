package encoder

import (
	"crypto/sha512"
	"encoding/base64"
	"io"
)

func HashAndEncode(input string) string {
	sha := sha512.New()
	io.WriteString(sha, input)
	return base64.StdEncoding.EncodeToString(sha.Sum(nil))
}
