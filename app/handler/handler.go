package handler

import (
	"fmt"
	"mxshs/redis-go/app/logger"
	"mxshs/redis-go/app/parser"
	"mxshs/redis-go/app/rdb"
	"mxshs/redis-go/app/replication"
	"mxshs/redis-go/app/store"
	"mxshs/redis-go/app/types"
	"mxshs/redis-go/app/utils"
	"net"
	"strconv"
	"strings"
	"time"
)

type Cmd int

const (
	_ Cmd = iota
	ECHO
)

const (
    BACKOFF_BASE = 2
    BACKOFF_EXP = 2
    BACKOFF_MAX = 4
)

type Handler struct {
    Listener net.Listener
    MasterConn net.Conn
    Slaves []net.Conn
	Info *replication.CoreInfo
	KV *store.KVStore
    RDB *rdb.RDBFile
    routes types.WorkerRouter
}

func NewSlaveHandler(la string, ma string) (*Handler, error) {
    master, err := net.Dial("tcp4", ma)
    if err != nil {
        return nil, err
    }

    l, err := net.Listen("tcp4", la)
    if err != nil {
        return nil, err
    }


    h := &Handler{
        Listener: l,
        MasterConn: master,
        Info: &replication.CoreInfo{
            Role: "slave",
        },
        KV: store.NewKVStore(),
        routes: types.WorkerRouter{},
    }


    err = h.syncToMaster()
    if err != nil {
        b := BACKOFF_BASE

        for ; err != nil; err = h.syncToMaster() {
            if b > BACKOFF_MAX {
                return nil, err
            }

            time.Sleep(time.Duration(b) * time.Second)
            b *= BACKOFF_EXP
        }
    }

    return h, nil
}

func NewMasterHandler(la string, kv *store.KVStore, db *rdb.RDBFile) (*Handler, error) {
    l, err := net.Listen("tcp4", la)
    if err != nil {
        return nil, err
    }


    h := &Handler{
        Listener: l,
        Slaves: make([]net.Conn, 0),
        Info: &replication.CoreInfo{
            Role: "master",
        },
        KV: store.NewKVStore(),
        RDB: rdb.DefaultRDB(),
        routes: types.WorkerRouter{},
    }

    return h, nil
}

func (h *Handler) Run() error {
    h.registerBaseWorkers()
    h.registerAuthWorkers()
    h.registerReplicationWorkers()
    h.registerKVWorkers()

    if h.Info.Role == "slave" {
		go h.handleConnection(h.MasterConn)
    }

	for {
		conn, err := h.Listener.Accept()
		if conn == nil || err != nil {
            logger.Logger.Debug(
                fmt.Sprintf("failed to open connection: %s", err.Error()),
            )
			continue
		}

        go h.handleConnection(conn)
	}
}

func (h *Handler) handleConnection(conn net.Conn) {
    defer conn.Close()

    rg := make([]byte, 1024)

    for {
        n, err := conn.Read(rg)
        if err != nil {
            logger.Logger.Debug(
                fmt.Sprintf("failed to read from connection with %s: %s", conn.RemoteAddr(), err.Error()),
            )
            return
        }

        if n > 0 {
            h.handleCommands(conn, rg[:n])
        }
    }
}

func (h *Handler) handleCommands(conn net.Conn, data []byte) {
    msgs, err := parser.Parse(string(data))
    if err != nil {
        logger.Logger.Info(err.Error())
        return
    }

    if len(msgs) == 0 {
        return
    }

	for _, msg := range msgs {
		args, ok := msg.Value.([]*types.Data)
        if !ok {
            continue
        }

        commandName := strings.ToLower(args[0].Value.(string))

        logger.Logger.Info(fmt.Sprintf("received %s command from %s", commandName, conn.RemoteAddr()))

	    worker, ok := h.routes[commandName]
        if !ok {
            continue
        }

        // check that command should be executed if slave or notify replicas if master
        if h.Info.Role == "slave" {
            if conn != h.MasterConn && !worker.ReplicaCallable {
                continue
            }
        } else if worker.CanReplicate {
            go func() {
                // TODO: use message size to index data slice, cuz now im wasting resources converting back to bytes
                toReplicate, _ := utils.Encode(msg)

                for _, slave := range h.Slaves {
                    logger.Logger.Info(fmt.Sprintf("propagating %s command to %s", commandName, slave.RemoteAddr())) 

                    _, err := slave.Write([]byte(*toReplicate))
                    if err != nil {
                        logger.Logger.Warn(
                            fmt.Sprintf("failed to propagate %s to %s: %s", commandName, slave.RemoteAddr(), err.Error()),
                        )
                        continue
                    }

                    logger.Logger.Info(fmt.Sprintf("successfully propagated %s command to %s", commandName, slave.RemoteAddr()))
                }
            }()
        }


        _, err := worker.Run(conn, msg)
        if err != nil {
            logger.Logger.Warn(
                fmt.Sprintf("failed to run command %s from %s: %s", commandName, conn.RemoteAddr(), err.Error()),
            )
            continue
        }

        logger.Logger.Info(fmt.Sprintf("successfully ran %s command from %s", commandName, conn.RemoteAddr()))

		if h.Info.Role == "slave" {
			h.Info.MasterReplOffset += msg.Sz
        }
    }
}

func (h *Handler) syncToMaster() error {
    conn := h.MasterConn
    if conn == nil {
        return fmt.Errorf("expected master connection to be open when creating a replica")
    }

    lp := strings.Split(h.Listener.Addr().String(), ":")[1]

	msg, err := utils.FastMessage(conn, "ping")
	if err != nil {
		return err
	}

	val, ok := msg[0].Value.(string)
	if !ok || val != "PONG" {
		return fmt.Errorf("received unexpected response from master at %s during handshake", conn.RemoteAddr().String())
	}

	msg, err = utils.FastMessage(conn, "REPLCONF", "listening-port", lp)
	if err != nil {
		return err
	}

	val, ok = msg[0].Value.(string)
	if !ok || val != "OK" {
		return fmt.Errorf("received unexpected response from master at %s during handshake", conn.RemoteAddr().String())
	}

	msg, err = utils.FastMessage(conn, "REPLCONF", "capa", "sync2")
	if err != nil {
		return err
	}

	val, ok = msg[0].Value.(string)
	if !ok || val != "OK" {
		return fmt.Errorf("received unexpected response from master at %s during handshake", conn.RemoteAddr().String())
	}

	msg, err = utils.FastMessage(conn, "PSYNC", "?", "-1")
	if err != nil {
		return err
	}

	val, ok = msg[0].Value.(string)
	if !ok {
		return fmt.Errorf("received unexpected response from master at %s during handshake", conn.RemoteAddr().String())
	}

	masterParams := strings.Split(val, " ")

	if len(masterParams) < 3 || masterParams[0] != "FULLRESYNC" {
		return fmt.Errorf("received unexpected response from master at %s during handshake", conn.RemoteAddr().String())
	}

    offset, err := strconv.Atoi(masterParams[2])
    if err != nil {
        return err
    }

	h.Info.MasterReplID = masterParams[1]
    h.Info.MasterReplOffset = offset
    
    temp := msg[1].Value.(string)

	curRDB, err := rdb.NewRDBFromString(temp)
	if err != nil {
		return err
	}

	h.RDB = curRDB

	return nil
}
