package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	Start(os.Args[1])
}

func Start(tcpAddrStr string)  {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpAddrStr)
	if err != nil {
		log.Printf("Resolve tcp addr failed :%v\n", err)
		return
	}

	//向服务器拨号
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Dial to server failed :%v \n", err)
		return
	}
	//向服务器发送消息
	go sendMsg(conn)

	//接收来自服务器的消息
	buff := make([]byte, 1024)
	for {
		length, err := conn.Read(buff)
		if err != nil {
			log.Printf("recv server msg failed %v \n", err)
			conn.Close()
			os.Exit(0)
			break
		}
		fmt.Println(string(buff[0:length]))
	}
}

func sendMsg(conn net.Conn)  {
	username := conn.LocalAddr().String()
	for {
		var input string
		//接收输入的消息 放到input变量中
		fmt.Scanln(&input)

		if input == "/q" || input == "/quit" {
			fmt.Println("Byebye ...")
			conn.Close()
			os.Exit(0)
		}

		if len(input) > 0 {
			msg := username + " say:" + input

			_, err := conn.Write([]byte(msg))
			if err != nil {
				conn.Close()
				break
			}
		}
	}
}