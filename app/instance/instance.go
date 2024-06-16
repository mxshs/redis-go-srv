package instance

import (
	"fmt"
	"mxshs/redis-go/app/handler"
)

func SpawnRedisInstance(port *int, ma, mp *string) error {
    var h *handler.Handler
    var err error

    //TODO: move to generic constructor (first i have to implement proper rdb handling)
    if ma != nil {
        h, err = handler.NewSlaveHandler(fmt.Sprintf("0.0.0.0:%d", *port), fmt.Sprintf("%s:%s", *ma, *mp))
        if err != nil {
            return err
        }

    } else {
        h, err = handler.NewMasterHandler(fmt.Sprintf("0.0.0.0:%d", *port), nil, nil)
        if err != nil {
            return err
        }
    }

    return h.Run()
}
