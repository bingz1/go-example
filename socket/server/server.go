package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	port := "9090"
	StartSer(port)
}

func StartSer(port string) {
	host := ":" + port

	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		log.Printf("resolve tcp addr failed:%v\n", err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("listen tcp port failed:%v\n", err)
		return
	}
	//建立连接池 用于广播消息
	conns := make(map[string]net.Conn)

	//消息通道
	messageChan := make(chan string, 10)

	//广播消息
	go BrodMessages(&conns, messageChan)

	//启动
	for {
		fmt.Printf("listening port %s ..\n", port)
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Accept failed:%v\n", err)
			continue
		}
		//把每个连接扔到连接池
		conns[conn.RemoteAddr().String()] = conn
		fmt.Println(conns)

		//处理消息
		go Handler(conn, &conns, messageChan)
	}

}

//向所有的连接发送广播消息
func BrodMessages(conns *map[string]net.Conn, messages chan string) {
	for {
		//不断从通道里读消息
		msg := <-messages
		fmt.Println(msg)

		for key, conn := range *conns {
			fmt.Println("connection is connected from ", key)
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Printf("broad message to %s failed:%v \n", key, err)
				delete(*conns, key)
			}
		}

	}
}

//处理客户端发送的消息 将其扔到通道里
func Handler(conn net.Conn, conns *map[string]net.Conn, messages chan string) {
	fmt.Println("connect from client", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Printf("read client message failed:%v\n", err)
			delete(*conns, conn.RemoteAddr().String())
			conn.Close()
			break
		}
		//把收到的消息写到通道中
		recvStr := string(buf[0:length])
		messages <- recvStr
	}
}
