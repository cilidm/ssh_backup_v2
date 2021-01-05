package main

import (
	"flag"
	"fmt"
	"log"
	"ssh_backup/app/router"
	"ssh_backup/app/service"
	"time"
)

var (
	runServer bool
	runMode   bool
)

func init() {
	flag.BoolVar(&runServer, "s", false, "开启http服务")
	flag.BoolVar(&runMode, "r", false, "开启传输任务")
	flag.Parse()

	service.CheckTargetDirs()
	service.InitLevelDB()
}

func main() {
	if !runMode && !runServer {
		log.Fatal("请至少选择一种模式运行")
	}

	if runMode {
		begin := time.Now()
		service.BeginRunMode()
		fmt.Println("传输结束，耗时", time.Since(begin))
	}

	if runServer {
		router.RunServer()
	}
}
