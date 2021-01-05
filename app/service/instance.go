package service

import (
	"github.com/cilidm/toolbox/levelDB"
	"github.com/cilidm/toolbox/logging"
	"ssh_backup/app/client"
	"ssh_backup/app/config"
	"ssh_backup/app/util"
)

// 运行模式包含ssh的先测试目标文件夹是否可以写入
func CheckTargetDirs() {
	if util.Find(config.Conf.RunMode, []int{Ltr, Rtr, Ltrto}) > 0 {
		_, err := ClientPathExists(config.Conf.Dir.Target, client.Instance())
		if err != nil {
			logging.Fatal("目标文件夹校验失败，请确定是否有写入权限", err.Error())
		}
	}
}

func InitLevelDB() {
	err := levelDB.InitServer("runtime")
	if err != nil {
		logging.Fatal(err.Error())
	}
}
