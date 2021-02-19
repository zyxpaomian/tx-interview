package main

import (
	"math/rand"
	go_http "net/http"
	"runtime"
	"time"
	cfg "tx-interview/common/configparse"
	log "tx-interview/common/formatlog"
	"tx-interview/common/mysql"
	"tx-interview/http"
	"tx-interview/http/handle"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	// 初始化配置 && 日志
	cfg.GlobalConf.CfgInit("./conf/simplehttp.ini")
	logname := cfg.GlobalConf.GetStr("common", "logname")
	loglevel := cfg.GlobalConf.GetStr("common", "loglevel")
	log.InitLog(logname, loglevel)

	// 初始化DB，如果DB连接不上则直接panic
	mysql.DB.InitConn()

	// 启动HTTP服务
	mux := http.New()
	handle.InitHandle(mux)
	srv := &go_http.Server{
		Handler:      mux.GetRouter(),
		Addr:         cfg.GlobalConf.GetStr("common", "httpsvr"),
		WriteTimeout: 15 * time.Hour,
		ReadTimeout:  15 * time.Hour,
	}
	log.Infof("[tx-interview] HTTP服务启动，监听地址为 %s", cfg.GlobalConf.GetStr("common", "httpsvr"))
	go srv.ListenAndServe()

	// 等待
	select {}
}
