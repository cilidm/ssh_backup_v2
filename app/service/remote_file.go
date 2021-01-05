package service

import (
	"fmt"
	"github.com/cilidm/toolbox/file"
	"github.com/cilidm/toolbox/levelDB"
	"github.com/cilidm/toolbox/logging"
	"github.com/cilidm/toolbox/str"
	"github.com/pkg/sftp"
	"path"
	"path/filepath"
	"ssh_backup/app/client"
	conf "ssh_backup/app/config"
	"ssh_backup/app/model"
	"ssh_backup/app/util"
	"time"
)

// begin from remote

// 遍历远端文件夹
func DirWalk(dirPath string, client *sftp.Client) {
	files, err := client.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		logging.Error(err)
		return
	}
	except := conf.Conf.FileInfo.ExceptDir
	for _, v := range files {
		if str.IsContain(except, v) {
			continue
		}
		stat, err := client.Stat(v)
		if err != nil {
			logging.Error(err)
			continue
		}
		if stat.IsDir() {
			DirWalk(v, client)
		} else {
			logging.Info("发现文件", v)
			pathHandle <- v
		}
	}
}

func PathMonitor() {
	for {
		val, ok := <-pathHandle
		if !ok {
			if len(maxReadChan) == 0 {
				done <- true
				break
			}
		} else {
			maxReadChan <- true
			go SyncSaveFileInfo(val)
		}
	}
}

// 将文件下载并存入ldb
func SyncSaveFileInfo(v string) {
	if checkFileStatus(v) {
		<-maxReadChan
		return
	}
	s, _ := client.Instance().Stat(v)
	var file model.FileInfo
	file.FileName = s.Name()
	file.FileSize = s.Size()
	file.Status = model.Prossing
	file.FileSource = v
	file.FileTarget = util.GetNewPath(conf.Conf.Dir.Source, conf.Conf.Dir.Target, v)
	file.CreatedAt = time.Now().Format(util.FormatTime)
	SaveFileFromLDBHandler(file)
	<-maxReadChan
}

func SaveFileFromLDBHandler(v model.FileInfo) {
	fmt.Println("开始传输文件", v.FileName)
	var mode string // 运行模式 目标本地/远端
	if util.Find(conf.Conf.RunMode, []int{Ltr, Rtr, Ltrto}) > 0 {
		mode = "ssh"
	} else if util.Find(conf.Conf.RunMode, []int{Rtl, Rtlto}) > 0 {
		mode = "local"
	}

	dir, fileName := path.Split(v.FileTarget)
	if mode == "ssh" { // 目标是云端
		if err := UploadFile(client.Instance(), client.TargetInstance(), dir, fileName, v.FileSource, v.FileSize); err != nil {
			v.Status = model.ProssErr
			levelDB.GetServer().Insert(util.GetLdbKey(conf.Conf.Dir.Source, v.FileSource), &v)
			logging.Error("upload file err :", err)
			return
		} else {
			AddLdbIndex(v)
			return
		}
	} else { // 目标是本地
		has, err := util.CheckFile(v.FileTarget) // 是否已存在
		if has != nil && err == nil {
			if has.Size() == v.FileSize {
				v.Status = model.Processed
				levelDB.GetServer().Insert(util.GetLdbKey(conf.Conf.Dir.Source, v.FileSource), &v)
				return
			}
		}
		file.IsNotExistMkDir(dir)
		if err := GetFile(client.Instance(), v.FileSource, v.FileTarget); err != nil {
			v.Status = model.ProssErr
			levelDB.GetServer().Insert(util.GetLdbKey(conf.Conf.Dir.Source, v.FileSource), &v)
			logging.Error(err)
			return
		}
		AddLdbIndex(v)
	}
}
