package handler

import (
	"fmt"
	"mxshs/redis-go/app/types"
	"mxshs/redis-go/app/utils"
	"net"
	"strconv"
	"time"
)

func (h *Handler) registerKVWorkers() {
    h.routes["set"] = types.CreateWorker(h.Set, true, false)

    h.routes["get"] = types.CreateWorker(h.Get, false, true)
}

// considers that the passed cmd array DOES INCLUDE the SET command itself
func (h *Handler) Set(conn net.Conn, cmd *types.Data) ([]byte, error) {
    if cmd.T != types.Array {
		return nil, fmt.Errorf("%v is not a valid command", *cmd)
    }

    args := cmd.Value.([]*types.Data)

    if len(args) < 3 {
		return nil, fmt.Errorf("%v is not a valid command", cmd)
    }

	key := args[1]
	if key.T != types.String {
		return nil, types.InvalidKey
	}

	k, ok := key.Value.(string)
	if !ok {
		return nil, types.InvalidKey
	}

	var expiry int64

    for idx := 3; idx < len(args); idx++ {
        opt := args[idx]
		if err := utils.ValidateBaseValue(opt, types.String); err != nil {
			return nil, err
		}

		switch opt.Value.(string) {
		case "px":
            if idx == len(args) - 1 {
                return nil, types.InvalidType
            }

            idx++
            rawExp := args[idx]

			if err := utils.ValidateBaseValue(rawExp, types.String); err != nil {
				return nil, err 
			}

			temp, err := strconv.Atoi(rawExp.Value.(string))
			if err != nil {
				return nil, types.InvalidType
			}

			expiry = time.Now().UnixMilli() + int64(temp)
		default:
			return nil, types.InvalidType
		}
	}

	err := h.KV.Set(k, args[2], int64(expiry))
	if err != nil {
		return nil, err
	}

    _, err = conn.Write([]byte(types.OK))
    if err != nil {
        return nil, err
    }

	return []byte(types.OK), nil
}

func (h *Handler) Get(conn net.Conn, cmd *types.Data) ([]byte, error) {
    if cmd.T != types.Array {
		return nil, fmt.Errorf("%v is not a valid command", *cmd)
    }

    args := cmd.Value.([]*types.Data)

    if len(args) < 2 {
		return nil, fmt.Errorf("%v is not a valid command", args)
    }

	if args[1].T != types.String {
		return nil, fmt.Errorf("%v is not a valid key", args[1].Value)
	}

	k, _ := args[1].Value.(string)

	value, err := h.KV.Get(k)
    if err != nil {
        return nil, err
    }

    result, err := utils.Encode(value)
    if err != nil {
        return nil, err
    }

    _, err = conn.Write([]byte(*result))
    if err != nil {
        return nil, err
    }

	return []byte(*result), nil
}
