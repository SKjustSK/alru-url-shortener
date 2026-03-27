package base62

import (
	"fmt"
	"os"

	"github.com/speps/go-hashids/v2"
)

func Encode(num uint64) string {
	hd := hashids.NewData()

	hd.Salt = os.Getenv("HASHID_SALT")

	hd.MinLength = 3

	h, err := hashids.NewWithData(hd)
	if err != nil {
		// Fallback to a string representation of the ID so the app doesn't crash
		return fmt.Sprintf("%d", num)
	}

	e, err := h.Encode([]int{int(num)})
	if err != nil {
		return fmt.Sprintf("%d", num)
	}

	return e
}
