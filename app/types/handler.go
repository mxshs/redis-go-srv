package types

import "net"

// Map for all commands available on handler
type WorkerRouter map[string]*Worker

// Type for a command handler (write a wrapper for handlers with non-matching signatures)
type Worker struct {
    Run Executor
    CanReplicate bool
    ReplicaCallable bool
}

type Executor func(net.Conn, *Data) ([]byte, error)

func CreateWorker(r Executor, cr, rc bool) *Worker {
    return &Worker{
        Run: r,
        CanReplicate: cr,
        ReplicaCallable: rc,
    }
}
