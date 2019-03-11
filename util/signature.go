package util

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
	"strings"
)

func CheckSignature(checkToken string, timestamp string, nonce string, signatureIn string) (result bool) {
	strlist := []string{checkToken, timestamp, nonce}
	sort.Strings(strlist)
	t := sha1.New()
	io.WriteString(t, strings.Join(strlist, ""))
	signatureGen := fmt.Sprintf("%x", t.Sum(nil))
	if signatureGen == signatureIn {
		result = true
	} else {
		result = false
	}
	return
}
