package main

import (
	"net"

	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:9090")
	if err != nil {
		panic(err)
	}
	var connlist []net.Conn
	for {
		oneconn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		connlist = append(connlist, oneconn)
		if len(connlist) == 2 {
			break
		}
	}
	//
	l.Close()
	//

}

func update(twoconn []net.Conn) {
	var turn int64
	for {
		// var
		var list []commmsg.JumpMSGOne
		for _, one := range twoconn { // 先读取两个客户端的指令
			msg := commmsg.ReadOnePack(one)
			list = append(list, msg.List...) // 拼成一个
		}
		//拼成一个
		var outmsg commmsg.JumpMSGTurn
		outmsg.List = list
		outmsg.Turn = turn
		turn++
		// 发出去
		commmsg.WriteJumpMSGTurn(twoconn, outmsg)
	}
}
