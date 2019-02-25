package impl

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	//存放websocket链接
	wsConn *websocket.Conn

	//用于存放客户端发来的数据
	inChan chan []byte
	//用于存放发给客户端的数据
	outChan chan []byte

	closeChan chan byte

	mutex sync.Mutex

	isClosed bool
}

//读取api
func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

//发送api
func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (conn *Connection) Close() {
	//线程安全的Close,可以并发多次调用 也叫做可重入的Close
	conn.wsConn.Close()

	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		close(conn.inChan)
		close(conn.outChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

func InitConnection(wsConn *websocket.Conn) (conn *Connection, err error) {
	conn = &Connection{
		wsConn:    wsConn,
		inChan:    make(chan []byte, 1000),
		outChan:   make(chan []byte, 1000),
		closeChan: make(chan byte, 1),
	}

	//启动读协程
	go conn.readLoop()
	//启动写协程
	go conn.writeLoop()

	return
}

//接收客户端发来的消息 放到队列
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}

		//容易阻塞到这里，等到inChain有空闲的位置
		select {
		case conn.inChan <- data: //把收到的数据先放到接收队列里
		case <-conn.closeChan:
			goto ERR
		}
	}

ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		select {
		case data = <-conn.outChan: //从出队列里取数据发给客户端
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}
