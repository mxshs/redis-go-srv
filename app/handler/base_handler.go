package handler

import (
	"fmt"
	"mxshs/redis-go/app/types"
	"mxshs/redis-go/app/utils"
	"net"
	"reflect"
	"strings"
)

func (h *Handler) registerBaseWorkers() {
    h.routes["echo"] = types.CreateWorker(h.runEcho, false, true)

    h.routes["info"] = types.CreateWorker(h.runInfo, false, true)

    h.routes["ping"] = types.CreateWorker(h.runPing, false, true)
}

func (h *Handler) runEcho(conn net.Conn, cmd *types.Data) ([]byte, error) {
    arr := cmd.Value.([]*types.Data)

    if len(arr) < 2 {
        return nil, fmt.Errorf("empty echo command: %v", arr)
    }

    response, err := utils.Encode(arr[1])
    if err != nil {
        return nil, err
    }

    _, err = conn.Write([]byte(*response))
    if err != nil {
        return nil, err
    }

    return []byte(*response), nil
}

func (h *Handler) runInfo(conn net.Conn, _ *types.Data) ([]byte, error) {
    info := types.Data{
        T: types.String,
    }

    replReflection := reflect.ValueOf(*h.Info)

    data := make([]string, replReflection.NumField())
    for i := 0; i < replReflection.NumField(); i++ {
        name := replReflection.Type().Field(i).Tag.Get("json")
        // sprint cuz it handles string representation for me
        data[i] = fmt.Sprintf("%s:%s\r\n", name, fmt.Sprint(replReflection.Field(i)))
    }

    info.Value = strings.Join(data, "")

    response, err := utils.Encode(&info)
    if err != nil {
        return nil, err
    }

    _, err = conn.Write([]byte(*response))
    if err != nil {
        return nil, err
    }

    return []byte(*response), nil
}

func (h *Handler) runPing(conn net.Conn, _ *types.Data) ([]byte, error) {
    _, err := conn.Write([]byte("+PONG\r\n"))
    if err != nil {
        return nil, err
    }

    return []byte("+PONG\r\n"), nil
}
