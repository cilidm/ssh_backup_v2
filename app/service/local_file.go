package service

import (
	"fmt"
	"os"
	"path/filepath"
	"ssh_backup/app/client"
	conf "ssh_backup/app/config"
	"ssh_backup/app/model"
	"ssh_backup/app/util"
	"strings"
	"time"

	"github.com/cilidm/toolbox/levelDB"
	"github.com/cilidm/toolbox/logging"
	"github.com/cilidm/toolbox/str"
)

// begin from local

func WalkPath(dir string) {
	paths, err := filepath.Glob(strings.TrimRight(dir, "/") + "/*")
	if err != nil {
		logging.Error(err)
		return
	}
	for _, v := range paths {
		if str.IsContain(conf.Conf.FileInfo.ExceptDir, v) || util.HasPathPrefix(conf.Conf.FileInfo.ExceptDir, v) {
			continue
		}
		stat, err := os.Stat(v)
		if err != nil {
			logging.Error(err)
			continue
		}
		if stat.IsDir() {
			WalkPath(v)
		} else {
			logging.Info("发现文件", v)
			pathHandle <- v
		}
	}
}

func LocalMonitor() {
	for {
		val, ok := <-pathHandle
		if !ok {
			if len(maxReadChan) == 0 {
				done <- true
				break
			}
		} else {
			maxReadChan <- true
			go GetFileInfo(val)
		}
	}
}

// 获取本地文件信息并开始传输
func GetFileInfo(v string) {
	if checkFileStatus(v) {
		<-maxReadChan
		return
	}
	s, err := os.Stat(v) // 此处可优化 只读一次
	if err != nil {
		logging.Error(err.Error())
		<-maxReadChan
		return
	}
	logging.Info("发现文件", s.Name())
	var file model.FileInfo
	file.FileName = s.Name()
	file.FileSize = s.Size()
	file.Status = model.Prossing
	file.FileSource = v
	file.FileTarget = util.GetNewPath(conf.Conf.Dir.Source, conf.Conf.Dir.Target, v)
	file.CreatedAt = time.Now().Format(util.FormatTime)
	UploadLocalFile(file)
	<-maxReadChan
}

func UploadLocalFile(v model.FileInfo) {
	fmt.Println("开始传输文件", v.FileName)
	logging.Info("开始传输文件", v.FileName)
	//dir, fileName := path.Split(v.FileTarget)
	dir, fileName := filepath.Split(v.FileTarget)
	// 源地址，
	if err := UploadFromLocal(client.Instance(), dir, fileName, v.FileSource, v.FileSize); err != nil {
		v.Status = model.ProssErr
		levelDB.GetServer().Insert(util.GetLdbKey(conf.Conf.Dir.Source, v.FileSource), &v)
		fmt.Println("upload file err :", err.Error())
		logging.Error("upload file err :", err)
		return
	}
	AddLdbIndex(v)
}
