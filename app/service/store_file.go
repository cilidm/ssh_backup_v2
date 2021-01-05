package service

import (
	"errors"
	"fmt"
	"github.com/cilidm/toolbox/file"
	"github.com/cilidm/toolbox/logging"
	"github.com/cilidm/toolbox/str"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"ssh_backup/app/config"
	"ssh_backup/app/model"
	"ssh_backup/app/store"
	"ssh_backup/app/util"
	"strconv"
	"strings"
	"time"
)

// 上传
func StoreMonitor() {
	for {
		val, ok := <-pathHandle
		if !ok {
			if len(maxReadChan) == 0 {
				done <- true
				break
			}
		} else {
			maxReadChan <- true
			go StartStore(config.Conf.Dir.OssPath, val)
		}
	}
}

func StartStore(dir, v string) {
	Upload2Store(dir, v)
	<-maxReadChan
}

func Upload2Store(dir, v string) {
	if checkFileStatus(v) {
		return
	}
	var conf store.Config
	str.CopyFields(&conf, config.Conf.Store[config.Conf.StoreType])
	conf.CloudType = store.GetStoreType(config.Conf.StoreType)
	cloud, err := store.NewCloudStore(conf, false)
	if err != nil {
		logging.Error(err.Error())
		return
	}
	stat, err := os.Stat(v)
	if err != nil {
		logging.Error(err.Error())
		return
	}
	filePath := v
	storePath := filepath.Join(dir, stat.Name())
	if cloud.IsExist(storePath) == nil {
		logging.Warn(v, "已存在")
		return
	}
	miMe := mime.TypeByExtension(path.Ext(filePath))
	if err := cloud.Upload(filePath, storePath, map[string]string{"Content-Type": miMe}); err != nil {
		fmt.Println(err)
	}
	AddLdbIndex(model.FileInfo{
		FileName:   stat.Name(),
		FileSize:   stat.Size(),
		FileSource: v,
		FileTarget: storePath,
		Status:     model.Processed,
		CreatedAt:  time.Now().Format(util.FormatTime),
	})
}

// 下载
func WalkStore() {
	var conf store.Config
	str.CopyFields(&conf, config.Conf.Store[config.Conf.StoreType])
	conf.CloudType = store.GetStoreType(config.Conf.StoreType)
	cloud, err := store.NewCloudStore(conf, false)
	if err != nil {
		logging.Error(err.Error())
		return
	}
	files, err := cloud.Lists(config.Conf.Dir.OssPath)
	if err != nil {
		logging.Error(err.Error())
		return
	}
	for _, v := range files {
		fmt.Println("发现文件", v.Name)
		link := cloud.GetSignURL(v.Name)
		pathHandle <- link
	}
}

func Store2LocalMonitor() {
	for {
		val, ok := <-pathHandle
		if !ok {
			if len(maxReadChan) == 0 {
				done <- true
				break
			}
		} else {
			maxReadChan <- true
			go saveStoreToLocal(val)
		}
	}
}

func saveStoreToLocal(link string) {
	beginSave(link)
	<-maxReadChan
}

var tempDir = "runtime/"

func beginSave(link string) {
	dir, name := getNameAndDir(link)
	fp := tempDir + dir + name
	err := file.IsNotExistMkDir(tempDir + dir)
	if err != nil {
		logging.Error(err.Error())
		return
	}
	err = DownloadFile(link, fp)
	if err != nil {
		logging.Error(err.Error())
		err = retryDown(link, fp) // 重试一次
		if err != nil {
			logging.Error(err.Error())
			return
		}
	}
	stat, err := os.Stat(fp)
	if err != nil {
		logging.Error(err.Error())
		return
	}
	AddLdbIndex(model.FileInfo{
		FileName:   stat.Name(),
		FileSize:   stat.Size(),
		FileSource: link,
		FileTarget: fp,
		Status:     model.Processed,
		CreatedAt:  time.Now().Format(util.FormatTime),
	})
}

func getNameAndDir(link string) (dir, name string) {
	shortLink := strings.Replace(link, config.Conf.Store[config.Conf.StoreType].PublicBucketDomain, "", 1)
	arr := strings.Split(shortLink, "/")
	if len(arr) == 1 {
		return "", arr[0]
	} else {
		return strings.TrimRight(shortLink, arr[len(arr)-1]), arr[len(arr)-1]
	}
}

func retryDown(url, path string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, res.Body)
	return err
}

func DownloadFile(url string, localPath string) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	tmpFilePath := localPath + ".download"
	//创建一个http client
	client := new(http.Client)
	client.Timeout = time.Second * 10 //设置超时时间
	//get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	//读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		logging.Error(url, "Get Content-Length Error")
		return err
	}
	if IsFileExist(localPath, fsize) {
		logging.Warn("IsFileExist")
		return err
	}
	//创建文件
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		logging.Error("body is null")
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	if err == nil {
		file.Close()
		err = os.Rename(tmpFilePath, localPath)
		if err != nil {
			logging.Error(err.Error())
			return err
		}
	}
	return err
}

func IsFileExist(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if filesize == info.Size() {
		logging.Warn("文件已存在！", info.Name(), info.Size(), info.ModTime())
		return true
	}
	del := os.Remove(filename)
	if del != nil {
		logging.Error("del err", del.Error())
	}
	return false
}
