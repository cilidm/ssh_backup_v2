package util

import (
	"crypto/md5"
	"fmt"
	"github.com/cilidm/toolbox/logging"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	FileIndexKey         = "file_index_"
	PicIndexKey          = "pic_index_"
	VideoIndexKey        = "video_index_"
	InitialID     uint64 = 100000001
)

var (
	PicExt   = []string{".jpg", ".jpeg", ".png", ".bmp", ".gif", ".webp"}
	VideoExt = []string{".mp4", ".wmv", ".avi", ".rm", ".mpeg", ".flv", ".3gp", ".mov"}
)

func GetPicIndexKey(id string) string {
	return PicIndexKey + id
}

func GetVideoIndexKey(id string) string {
	return VideoIndexKey + id
}

func GetFileIndexKey(id string) string {
	return FileIndexKey + id
}

func GetLdbKey(path, md string) string {
	return "LDB_" + path + "_" + md
}

func GetLdbPreKey(path string) string {
	return "LDB_" + path
}

// 原文件夹地址/home 目标文件夹地址/newPath 文件地址/home/a/b/ca.txt 则返回/newPath/a/b/ca.txt
func GetNewPath(source, target, dst string) string {
	return filepath.Join(target, strings.ReplaceAll(dst, source, ""))
}

func CheckFile(path string) (os.FileInfo, error) {
	file, err := os.Stat(path)
	if err == nil {
		return file, nil
	}
	if os.IsNotExist(err) {
		return nil, nil
	}
	return nil, err
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetMdByPath(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		logging.Error(err.Error())
		return "", err
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		logging.Error(err.Error())
		return "", err
	}
	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}


// 检测并创建文件夹
//func PathExists(path string) (bool, error) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true, nil
//	}
//	if os.IsNotExist(err) {
//		os.MkdirAll(path, os.ModePerm)
//		return false, nil
//	}
//	return false, err
//}

//func CheckErr(err error, msg string) {
//	if err != nil {
//		log.Printf(msg)
//		os.Exit(500)
//	}
//}

//func SizeFormat(size float64) string {
//	units := []string{"Byte", "KB", "MB", "GB", "TB"}
//	n := 0
//	for size > 1024 {
//		size /= 1024
//		n += 1
//	}
//
//	return fmt.Sprintf("%.2f %s", size, units[n])
//}

//func CopyFields(a interface{}, b interface{}, fields ...string) (err error) {
//	at := reflect.TypeOf(a)
//	av := reflect.ValueOf(a)
//	bt := reflect.TypeOf(b)
//	bv := reflect.ValueOf(b)
//	// 简单判断下
//	if at.Kind() != reflect.Ptr {
//		err = fmt.Errorf("a must be a struct pointer")
//		return
//	}
//	av = reflect.ValueOf(av.Interface())
//	// 要复制哪些字段
//	_fields := make([]string, 0)
//	if len(fields) > 0 {
//		_fields = fields
//	} else {
//		for i := 0; i < bv.NumField(); i++ {
//			_fields = append(_fields, bt.Field(i).Name)
//		}
//	}
//	if len(_fields) == 0 {
//		fmt.Println("no fields to copy")
//		return
//	}
//	// 复制
//	for i := 0; i < len(_fields); i++ {
//		name := _fields[i]
//		f := av.Elem().FieldByName(name)
//		bValue := bv.FieldByName(name)
//		// a中有同名的字段并且类型一致才复制
//		if f.IsValid() && f.Kind() == bValue.Kind() {
//			f.Set(bValue)
//		} else {
//			//fmt.Printf("no such field or different kind, fieldName: %s\n", name)
//		}
//	}
//	return
//}

// 判断字符串是否在数组里
//func IsContain(items []string, item string) bool {
//	for _, eachItem := range items {
//		if eachItem == item {
//			return true
//		}
//	}
//	return false
//}
