package commmsg

import (
	"encoding/binary"
	"net"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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
	if err := ReadFixBytes(conn, buffer); err != nil {
		panic(err)
	}
	datalen := binary.BigEndian.Uint64(buffer)
	return int64(datalen)
}

func ReadOnePack(conn net.Conn, vaddr protoreflect.ProtoMessage) {
	datalen := ReadBytesToInt(conn)
	// fmt.Println("datalen := ReadBytesToInt(conn)", datalen)
	buffer := make([]byte, datalen)
	ReadFixBytes(conn, buffer)
	proto.Unmarshal(buffer, vaddr)
}

func WriteJumpMSGTurn(connlist []net.Conn, vaddr protoreflect.ProtoMessage) {
	data, err := proto.Marshal(vaddr)
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
