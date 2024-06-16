package handler

import (
	"mxshs/redis-go/app/types"
	"net"
)

func (h *Handler) registerAuthWorkers() {
    h.routes["hello"] = types.CreateWorker(h.runHello, false, true)
}

// TODO: hello cmd for redis client auth, gotta add some fields to handler
func (h *Handler) runHello(conn net.Conn, cmd *types.Data) ([]byte, error) {
    response := []byte(
        "%3\r\n+server\r\n+redis\r\n+version\r\n:123\r\n+proto\r\n:3\r\n",
    )

    _, err := conn.Write(response)
    if err != nil {
        return nil, err
    }

    return response, nil
}
