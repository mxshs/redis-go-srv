package rdb

import "encoding/hex"

const EMPTYRDB string = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"

type RDBFile struct {
    RDB []byte
    Location string
}

func NewRDBFromHex(hx []byte) (*RDBFile, error) {
	body := make([]byte, hex.DecodedLen(len(hx)))

	_, err := hex.Decode(body, hx)
	if err != nil {
		return nil, err
	}

	return &RDBFile{
		RDB:      body,
        Location: "",
	}, nil
}

func NewRDBFromString(rdb string) (*RDBFile, error) {
	return &RDBFile{
		RDB:      []byte(rdb),
		Location: "",
	}, nil
}

func DefaultRDB() *RDBFile {
	rdb, _ := NewRDBFromHex([]byte(EMPTYRDB))

	return rdb
}
