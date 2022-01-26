package commmsg

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

type JumpMSGTurn struct {
	Turn int64
	List []JumpMSGOne
}

type JumpMSGOne struct {
	Kind     string // move jump choose login
	UID      string
	Which    int64 // 哪一个
	MoveValX int64
	MoveValZ int64
}

func ReadFixBytes(conn net.Conn, buffer []byte) error {
	already := 0
	for {
		n, err := conn.Read(buffer[already:])
		if err != nil {
			return err
		}
		already += n
		if already == len(buffer) {
			break
		}
	}
	return nil
}
func WriteFixedBytes(conn net.Conn, buffer []byte) error {
	already := 0
	for {
		n, err := conn.Write(buffer[already:])
		if err != nil {
			return err
		}
		already += n
		if already == len(buffer) {
			break
		}
	}
	return nil
}

func ReadBytesToInt(conn net.Conn) int64 {
	buffer := make([]byte, 8)
	ReadFixBytes(conn, buffer)
	datalen := binary.BigEndian.Uint64(buffer)
	return int64(datalen)
}

func ReadOnePack(conn net.Conn) JumpMSGTurn {
	datalen := ReadBytesToInt(conn)
	buffer := make([]byte, datalen)
	ReadFixBytes(conn, buffer)
	var res JumpMSGTurn
	json.Unmarshal(buffer, &res)
	return res
}

func WriteJumpMSGTurn(connlist []net.Conn, msg JumpMSGTurn) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	arr := make([]byte, 8)
	binary.BigEndian.PutUint64(arr, uint64(len(data)))
	for _, oneconn := range connlist {
		WriteFixedBytes(oneconn, arr)
		WriteFixedBytes(oneconn, data)
	}
}
