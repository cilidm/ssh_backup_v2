package test

import (
	"fmt"
	"github.com/cilidm/toolbox/str"
	"log"
	"mime"
	"path"
	"path/filepath"
	"ssh_backup/app/config"
	"ssh_backup/app/store"
	"testing"
)

func TestStore(t *testing.T) {
	var conf store.Config
	str.CopyFields(&conf,config.Conf.Store[config.Conf.StoreType])
	conf.CloudType = store.GetStoreType(config.Conf.StoreType)
	cloud,err := store.NewCloudStore(conf,false)
	if err != nil{
		log.Fatal(err.Error())
	}
	filePath := "views/index.html"
	storePath := filepath.Join("test",filePath)
	miMe := mime.TypeByExtension(path.Ext(filePath))
	if err := cloud.Upload(filePath,storePath,map[string]string{"Content-Type": miMe});err != nil{
		fmt.Println(err)
	}

	files,err := cloud.Lists("")		// 腾讯云暂时没有lists
	if err != nil{
		log.Fatal(err.Error())
	}
	for _,v := range files{
		fmt.Println(v.Name)
	}
}
