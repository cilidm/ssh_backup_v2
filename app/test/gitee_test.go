package test

import (
	"fmt"
	"log"
	"ssh_backup/app/store"
	"testing"
	"time"
)

// test gitee
// 递归会403
func getContents(p string) {
	files, err := store.GetContents(p)
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, v := range files {
		if v.Type == "dir" {
			time.Sleep(time.Second * 2)
			getContents(p)
		} else {
			if v.Name == ".gitignore" {
				continue
			}
			fmt.Println("begin", v.Name)
			blob, err := store.GetBlob(v.Sha)
			if err != nil {
				continue
			}
			store.Base64ToFile(blob.Content, "runtime/gitee/"+v.Name)
			fmt.Println(v.Name, "save success")
			time.Sleep(time.Second * 2)
		}
	}
}

func TestGitee(t *testing.T) {
	err := store.CreateFile("runtime/logs/log_20201225.log", "log/1.log", "test")
	fmt.Println(err)
}
