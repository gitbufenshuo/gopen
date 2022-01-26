package main

import (
	"fmt"
	"net"
	"os"

	"github.com/gitbufenshuo/gopen/example/jumpjump/commmsg"
)

func main() {
	addr := os.Args[1]
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
		if len(connlist) == 2 {
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
		var list []commmsg.JumpMSGOne
		for _, one := range twoconn { // 先读取两个客户端的指令
			msg := commmsg.ReadOnePack(one)
			list = append(list, msg.List...) // 拼成一个
			fmt.Printf("[%d]Read %s\n", turn, one.RemoteAddr().String())
		}
		//拼成一个
		var outmsg commmsg.JumpMSGTurn
		outmsg.List = list
		outmsg.Turn = turn
		// 发出去
		commmsg.WriteJumpMSGTurn(twoconn, outmsg)
		fmt.Println("send", turn)
		turn++
	}
}
