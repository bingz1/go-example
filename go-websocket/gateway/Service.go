package gateway

import (
	"net/http"
	"time"
)

func handlePushAll(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello World"))
}

func InitService() (err error) {

	var (
		mux    *http.ServeMux
		server *http.Server
	)

	//路由
	mux = http.NewServeMux()
	mux.HandleFunc("/push/all", handlePushAll)

	server = &http.Server{
		Addr:         ":7799",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      mux,
	}

	// 拉起服务
	go server.ListenAndServe()
	return
}
