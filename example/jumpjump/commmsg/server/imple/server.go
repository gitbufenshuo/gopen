package imple

import (
	"fmt"
	"net"

	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg"
	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg/protodefine/pgocode/jump"
)

func Main(count int, addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("begin listen", addr)
	var connlist []net.Conn
	for {
		oneconn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println(oneconn.RemoteAddr(), "In")
		connlist = append(connlist, oneconn)
		if len(connlist) == count {
			break
		}
	}
	//
	l.Close()
	//
	update(connlist)
}

func update(twoconn []net.Conn) {
	var turn int64
	for {
		// var
		var totalmsg jump.JumpMSGTurn
		for _, one := range twoconn { // 先读取两个客户端的指令
			var onemsg jump.JumpMSGTurn
			commmsg.ReadOnePack(one, &onemsg)
			totalmsg.List = append(totalmsg.List, onemsg.List...)
			//		fmt.Printf("[%d]Read %s %v\n", turn, one.RemoteAddr().String(), onemsg.List)
		}
		//拼成一个
		totalmsg.Turn = turn
		// 发出去
		commmsg.WriteJumpMSGTurn(twoconn, &totalmsg)
		//	fmt.Println("send", turn, totalmsg.List)
		turn++
	}
}
