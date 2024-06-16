package handler

import (
	"fmt"
	"mxshs/redis-go/app/logger"
	"mxshs/redis-go/app/types"
	"net"
	"strconv"
)

func (h *Handler) registerReplicationWorkers() {
    h.routes["replconf"] = types.CreateWorker(h.runReplconf, false, false)

    h.routes["psync"] = types.CreateWorker(h.runPsync, false, false)
}

func (h *Handler) runReplconf(conn net.Conn, cmd *types.Data) ([]byte, error) {
    args := cmd.Value.([]*types.Data)
    if len(args) < 3 {
        return nil, fmt.Errorf("malformed replconf message: %v", *cmd)
    }
    
    var result []byte

    if h.Info.Role == "slave" {
        logger.Logger.Info(fmt.Sprintf("%s attempted to sync this replica", conn.LocalAddr()))

        if conn != h.MasterConn {
            return nil, fmt.Errorf("client connection attempted to replconf a replica")
        }

        offset := strconv.Itoa(h.Info.MasterReplOffset)

        result = []byte(
            fmt.Sprintf(
                "*3\r\n$8\r\nREPLCONF\r\n$3\r\nACK\r\n$%d\r\n%s\r\n",
                len(offset), 
                offset,
            ),
        )
    } else {
        logger.Logger.Info(fmt.Sprintf("replica at %s acknowledged synchronization", conn.LocalAddr()))

        result = []byte(types.OK)
    }

    _, err := conn.Write([]byte(result))
    if err != nil {
        return nil, err
    }

    return result, nil
}

// TODO: Partial sync (for backlog update)
func (h *Handler) runPsync(conn net.Conn, cmd *types.Data) ([]byte, error) {
    h.Slaves = append(h.Slaves, conn)

    result := []byte(
        fmt.Sprintf(
            "+FULLRESYNC %s 0\r\n$%d\r\n%s",
            h.Info.MasterReplID,
            len(h.RDB.RDB),
            h.RDB.RDB,
        ),
    )

    _, err := conn.Write([]byte(result))
    if err != nil {
        return nil, err
    }

    return result, nil
}
