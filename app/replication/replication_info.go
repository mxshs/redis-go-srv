package replication

import "net"

type CoreInfo struct {
	Role                       string `json:"role"`
	MasterReplID               string `json:"master_replid"`
	ConnectedSlaves            int
	MasterReplOffset           int `json:"master_repl_offset"`
	SecondReplOffset           int
	ReplBacklogActive          int
	ReplBacklogSize            int
	ReplBacklogFirstByteOffset int
	ReplBacklogHistLen         int
}

type ReplicationInfo struct {
    Info CoreInfo
    // Used by master to propagate messages to replicas
	Slaves                     []net.Conn
    // Used by slave instances to handle messages from master
	MasterConn                 net.Conn
}

func NewReplicationInfo(role string) *CoreInfo {
    return &CoreInfo{
        Role: role,
    }
}

