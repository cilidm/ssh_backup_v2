package main

import (
	"fmt"
	"github.com/cilidm/toolbox/str"
	"log"
	"mime"
	"path"
	"path/filepath"
	"ssh_backup/app/config"
	"ssh_backup/app/service"
	"ssh_backup/app/store"
	"strings"
)

func TestTrim() {
	a := "bbb/aaa/"
	fmt.Println(strings.TrimRight(a, "/"))
}

func TestStore() {
	var conf store.Config
	str.CopyFields(&conf, config.Conf.Store[config.Conf.StoreType])
	conf.CloudType = store.GetStoreType(config.Conf.StoreType)
	cloud, err := store.NewCloudStore(conf, false)
	if err != nil {
		log.Fatal(err.Error())
	}
	filePath := "views/index.html"
	storePath := filepath.Join("test", filePath)
	miMe := mime.TypeByExtension(path.Ext(filePath))
	if err := cloud.Upload(filePath, storePath, map[string]string{"Content-Type": miMe}); err != nil {
		fmt.Println(err)
	}
}

func TestStoreList() {
	var conf store.Config
	str.CopyFields(&conf, config.Conf.Store[config.Conf.StoreType])
	conf.CloudType = store.GetStoreType(config.Conf.StoreType)
	cloud, err := store.NewCloudStore(conf, false)
	if err != nil {
		log.Fatal(err.Error())
	}
	files, err := cloud.Lists("") // 腾讯云暂时没有lists
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, v := range files {
		fmt.Println(v.Name)
	}
}

func getStoreDir()  {
	fmt.Println(config.Conf.Store[config.Conf.StoreType].PublicBucketDomain)
}

func main() {
	service.DownloadFile("http://alinkimg.mooncake.pw/qiniu_util/log_20210104.log","runtime/log_20210104.log")

	//TestStoreList()
	//getStoreDir()
}
