package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"sync"
)

type TomlConfig struct {
	RunMode   int `toml:"run_mode" validate:"required,min=1,max=10"`
	Dir       dir
	FileInfo  fileinfo
	Ssh       ssh
	StoreType string `toml:"store_type" validate:"omitempty,oneof=minio oss cos qiniu upyun obs bos gitee"`
	Store     map[string]store
	Gitee     gitee
}

type dir struct {
	Source  string `toml:"source" validate:"required"`
	Target  string `toml:"target" validate:"required"`
	OssPath string `toml:"oss_path"`
}

type fileinfo struct {
	ExceptDir  []string `toml:"except_dir"`
	MaxChannel int      `toml:"max_channel" validate:"min=1,max=20"`
}

type ssh struct {
	SourceHost string `toml:"source_host"`
	SourceUser string `toml:"source_user"`
	SourcePwd  string `toml:"source_pwd"`
	SourcePort int    `toml:"source_port"`
	TargetHost string `toml:"target_host"`
	TargetUser string `toml:"target_user"`
	TargetPwd  string `toml:"target_pwd"`
	TargetPort int    `toml:"target_port"`
}

type store struct {
	AccessKey           string `toml:"access_key"`
	SecretKey           string `toml:"secret_key"`
	Endpoint            string `toml:"endpoint"`
	PublicBucket        string `toml:"public_bucket"`
	PublicBucketDomain  string `toml:"public_bucket_domain"`
	PrivateBucket       string `toml:"private_bucket"`
	PrivateBucketDomain string `toml:"private_bucket_domain"`
	Expire              string `toml:"expire"`
	Region              string `toml:"region"`
	AppId               string `toml:"app_id"`
}

type gitee struct {
	Token string
	Owner string
	Repo  string
}

var (
	Conf TomlConfig
	once sync.Once
)

func init() {
	once.Do(func() {
		_, err := toml.DecodeFile("./conf.toml", &Conf)
		if err != nil {
			log.Fatal(err.Error())
		}
		validate := validator.New()
		err = validate.Struct(Conf)
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println("config配置文件", err.Field(), "字段配置错误")
			}
			os.Exit(-1)
		}
	})
}
