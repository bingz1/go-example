package gateway

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wantgo/go-example/go-websocket/impl"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{

		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		conn   *impl.Connection
		err    error
		data   []byte
	)

	//完成应答 在header中存放如下参数
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	//conn 客户端发起的连接
	if conn, err = impl.InitConnection(wsConn); err != nil {
		goto ERR
	}

	//发送心跳 维持连接
	go func() {
		var (
			err error
		)
		for {
			//每隔一秒发送一次心跳
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			fmt.Println(err)
			//收到消息 处理业务逻辑  处理完向相对应的客户端发送数据
			goto ERR
		}

		if err = conn.WriteMessage(data); err != nil {
			fmt.Println(err)
			goto ERR
		}
	}

ERR:
	conn.Close()
}

//开启 websocket 服务
func InitWSServer() (err error) {

	http.HandleFunc("/ws", wsHandler)

	url := "0.0.0.0:" + fmt.Sprintf("%d", G_config.WsPort)

	go http.ListenAndServe(url, nil)

	return
}
