package utils

import (
	"mxshs/redis-go/app/parser"
	"mxshs/redis-go/app/types"
	"net"
)

func FastMessage(conn net.Conn, message ...any) ([]*types.Data, error) {
	buf := make([]byte, 1024)

	msg, err := Encode(types.FastCommand(message))
	if err != nil {
		return nil, err
	}

	_, err = conn.Write([]byte(*msg))
	if err != nil {
		return nil, err
	}

	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	response, err := parser.Parse(string(buf[:n]))
	if err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, types.EmptyMessage
	}

	return response, nil
}
