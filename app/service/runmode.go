package service

import (
	"fmt"
	"github.com/cilidm/toolbox/logging"
	"os"
	"ssh_backup/app/client"
	"ssh_backup/app/config"
)

const (
	Ltr   = 1
	Lto   = 2
	Rto   = 3
	Rtr   = 4
	Oto   = 5
	Rtl   = 6
	Otl   = 7
	Ltrto = 8
	Rtlto = 9
	Otlts = 10
)

type runMode func()

func BeginRunMode() {
	RunMode := make(map[int]runMode)
	RunMode[Ltr] = Local2Romote
	RunMode[Lto] = Local2Oss
	RunMode[Rto] = Remote2Oss
	RunMode[Rtr] = Remote2Remote
	RunMode[Oto] = Oss2Oss
	RunMode[Rtl] = Remote2Local
	RunMode[Otl] = Oss2Local
	RunMode[Ltrto] = Local2Remote2Oss
	RunMode[Rtlto] = Remote2Local2Oss
	RunMode[Otlts] = Oss2Local2Remote
	RunMode[config.Conf.RunMode]()
}

func Local2Romote() {
	go LocalMonitor()
	startFromLocal()
}

func Local2Oss() {
	go StoreMonitor()
	startFromLocal()
}

func Remote2Oss() {
}

func Remote2Remote() {
	logging.Info("begin remote to remote")
	go PathMonitor() // 监听dir_walk
	startFromRemote()
}

func Oss2Oss() {
}

func Remote2Local() {
	logging.Info("begin remote to local")
	go PathMonitor() // 监听dir_walk
	startFromRemote()
}

func Oss2Local() {
	logging.Info("begin store to local")
	go Store2LocalMonitor()
	startFromStore()
}

func Local2Remote2Oss() {
}

func Remote2Local2Oss() {
}

func Oss2Local2Remote() {

}

func startFromStore()  {
	WalkStore()
	close(pathHandle)
	select {
	case <-done:
		logging.Info("任务结束")
	}
}

func startFromLocal() {
	WalkPath(config.Conf.Dir.Source)
	close(pathHandle)

	select {
	case <-done:
		logging.Info("任务结束")
	}
}

func startFromRemote() {
	cli := client.Instance()
	if cli == nil {
		fmt.Println("无法连接远程服务器")
		os.Exit(1)
	}
	DirWalk(config.Conf.Dir.Source, cli)
	close(pathHandle) // dir_walk结束，关闭path_chan

	// 阻塞，等待传输结束
	select {
	case <-done:
		logging.Info("任务结束")
	}
}
