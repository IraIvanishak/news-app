package domain

import (
	"encoding/hex"
	"errors"
)

type ID [12]byte

func IDFromHex(s string) (ID, error) {
	if len(s) != 24 {
		return ID{}, errors.New("invalid hex length")
	}

	var oid [12]byte
	_, err := hex.Decode(oid[:], []byte(s))
	if err != nil {
		return ID{}, err
	}

	return oid, nil
}
