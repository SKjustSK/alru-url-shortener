package base62

import "strings"

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base = uint64(len(alphabet))

func Encode(num uint64) string {
	if num == 0 {
		return string(alphabet[0])
	}

	var b strings.Builder
	// Pre-allocating 12 bytes covers the max value of a uint64 in Base62.
	// This prevents the builder from having to resize memory while looping.
	b.Grow(12)

	for num > 0 {
		b.WriteByte(alphabet[num%base])
		num /= base
	}

	return reverse(b.String())
}

func reverse(s string) string {
	bytes := []byte(s)

	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return string(bytes)
}
