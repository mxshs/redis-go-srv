package main

import (
	"flag"
	"mxshs/redis-go/app/instance"
	"mxshs/redis-go/app/logger"
)

func main() {
	port := flag.Int("port", 6379, "port for the current redis instance")
	replicaOf := flag.String("replicaof", "", "port for the current redis instance")

	flag.Parse()

    if *replicaOf != "" {
        logger.Setup(true, *replicaOf)
        err := instance.SpawnRedisInstance(port, replicaOf, &flag.Args()[0])
        if err != nil {
            panic(err)
        }
    } else {
        logger.Setup(false, "")
        err := instance.SpawnRedisInstance(port, nil, nil)
        if err != nil {
            panic(err)
        }
    }
}
