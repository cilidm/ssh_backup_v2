package store

import (
	"fmt"
	"github.com/cilidm/toolbox/store/cloud"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	AccessKey           string
	SecretKey           string
	Endpoint            string
	Region              string
	AppId               string
	PublicBucket        string
	PublicBucketDomain  string
	PrivateBucket       string
	PrivateBucketDomain string
	Expire              int64
	CloudType           StoreType
}

type StoreType string

var (
	Oss   StoreType = "oss"   //oss存储
	Minio StoreType = "minio" //minio存储
	Cos   StoreType = "cos"   //腾讯云存储
	Obs   StoreType = "obs"   //华为云存储
	Bos   StoreType = "bos"   //百度云存储
	Qiniu StoreType = "qiniu" //七牛云储存
	Upyun StoreType = "upyun" //又拍云存储
)

func GetStoreType(st string) StoreType {
	store := make(map[string]StoreType)
	store["oss"] = Oss
	store["minio"] = Minio
	store["cos"] = Cos
	store["obs"] = Obs
	store["bos"] = Bos
	store["qiniu"] = Qiniu
	store["upyun"] = Upyun
	return store[st]
}

type ConfigOss struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigMinio struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigCos struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Region              string `store:"region" fieldName:"Region"`
	AppId               string `store:"app-id" fieldName:"AppId"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigBos struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigObs struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigQiniu struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigUpYun struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type CloudStore struct {
	Private       bool
	StoreType     StoreType
	CanGZIP       bool
	Client        interface{}
	Config        interface{}
	Expire        int64
	PublicDomain  string
	PrivateDomain string
}

func (c *CloudStore) Lists(object string) (files []cloud.File, err error) {
	switch c.StoreType {
	case Cos:
		files, err = c.Client.(*cloud.COS).Lists(object)
	case Oss:
		files, err = c.Client.(*cloud.OSS).Lists(object)
	case Bos:
		files, err = c.Client.(*cloud.BOS).Lists(object)
	case Obs:
		files, err = c.Client.(*cloud.OBS).Lists(object)
	case Upyun:
		files, err = c.Client.(*cloud.UpYun).Lists(object)
	case Minio:
		files, err = c.Client.(*cloud.MinIO).Lists(object)
	case Qiniu:
		files, err = c.Client.(*cloud.QINIU).Lists(object)
	}
	return
}

func (c *CloudStore) Upload(tmpFile, saveFile string, headers ...map[string]string) (err error) {
	switch c.StoreType {
	case Cos:
		err = c.Client.(*cloud.COS).Upload(tmpFile, saveFile, headers...)
	case Oss:
		err = c.Client.(*cloud.OSS).Upload(tmpFile, saveFile, headers...)
	case Bos:
		err = c.Client.(*cloud.BOS).Upload(tmpFile, saveFile, headers...)
	case Obs:
		err = c.Client.(*cloud.OBS).Upload(tmpFile, saveFile, headers...)
	case Upyun:
		err = c.Client.(*cloud.UpYun).Upload(tmpFile, saveFile, headers...)
	case Minio:
		err = c.Client.(*cloud.MinIO).Upload(tmpFile, saveFile, headers...)
	case Qiniu:
		err = c.Client.(*cloud.QINIU).Upload(tmpFile, saveFile, headers...)
	}
	return
}

func (c *CloudStore) Delete(objects ...string) (err error) {
	switch c.StoreType {
	case Cos:
		err = c.Client.(*cloud.COS).Delete(objects...)
	case Oss:
		err = c.Client.(*cloud.OSS).Delete(objects...)
	case Bos:
		err = c.Client.(*cloud.BOS).Delete(objects...)
	case Obs:
		err = c.Client.(*cloud.OBS).Delete(objects...)
	case Upyun:
		err = c.Client.(*cloud.UpYun).Delete(objects...)
	case Minio:
		err = c.Client.(*cloud.MinIO).Delete(objects...)
	case Qiniu:
		err = c.Client.(*cloud.QINIU).Delete(objects...)
	}
	return
}

// err 返回 nil，表示文件存在，否则表示文件不存在
func (c *CloudStore) IsExist(object string) (err error) {
	switch c.StoreType {
	case Cos:
		err = c.Client.(*cloud.COS).IsExist(object)
	case Oss:
		err = c.Client.(*cloud.OSS).IsExist(object)
	case Bos:
		err = c.Client.(*cloud.BOS).IsExist(object)
	case Obs:
		err = c.Client.(*cloud.OBS).IsExist(object)
	case Upyun:
		err = c.Client.(*cloud.UpYun).IsExist(object)
	case Minio:
		err = c.Client.(*cloud.MinIO).IsExist(object)
	case Qiniu:
		err = c.Client.(*cloud.QINIU).IsExist(object)
	}
	return
}

func (c *CloudStore) PingTest() (err error) {
	tmpFile := "mybed-test-file.txt"
	saveFile := "mybed-test-file.txt"
	text := "hello world"

	defer func() {
		if err != nil {
			err = fmt.Errorf("Bucket是否私有：%v，错误信息：%v", c.Private, err.Error())
		}
	}()

	err = ioutil.WriteFile(tmpFile, []byte(text), os.ModePerm)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	if err = c.Upload(tmpFile, saveFile); err != nil {
		return
	}
	if err = c.IsExist(saveFile); err != nil {
		return
	}
	if err = c.Delete(saveFile); err != nil {
		return
	}
	return
}

//设置默认图片
//@param                picture             图片文件或者图片文件md5等
//@param                ext                 图片扩展名，如果图片文件参数(picture)的值为md5时，需要加上后缀扩展名
//@return               link                图片url链接
func (c *CloudStore) getImageFromCloudStore(picture string, ext ...string) (link string) {
	if len(ext) > 0 {
		picture = picture + "." + ext[0]
	} else if !strings.Contains(picture, ".") && len(picture) > 0 {
		picture = picture + ".jpg"
	}
	if c == nil || c.Client == nil {
		return
	}

	return c.GetSignURL(picture)
}

func (c *CloudStore) GetSignURL(object string) (link string) {
	var err error
	switch c.StoreType {
	case Cos:
		link, err = c.Client.(*cloud.COS).GetSignURL(object, c.Expire)
	case Oss:
		link, err = c.Client.(*cloud.OSS).GetSignURL(object, c.Expire)
	case Bos:
		link, err = c.Client.(*cloud.BOS).GetSignURL(object, c.Expire)
	case Obs:
		link, err = c.Client.(*cloud.OBS).GetSignURL(object, c.Expire)
	case Upyun:
		link, err = c.Client.(*cloud.UpYun).GetSignURL(object, c.Expire)
	case Minio:
		link, err = c.Client.(*cloud.MinIO).GetSignURL(object, c.Expire)
	case Qiniu:
		link, err = c.Client.(*cloud.QINIU).GetSignURL(object, c.Expire)
	}
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (c *CloudStore) GetPublicDomain() (domain string) {
	object := "test.test"
	link := c.GetSignURL(object)
	return strings.TrimRight(strings.Split(link, object)[0], "/")
}
