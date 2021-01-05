package store

import (
	"errors"
	"github.com/cilidm/toolbox/store/cloud"
	"reflect"
)

// 创建云存储
// config 相应存储的struct指针
// storetype 云存储对应的名字
// private false不使用普通网络
func NewCloudStore(storeConf Config,private bool) (cs *CloudStore, err error) {
	storeType := storeConf.CloudType
	config := GetConfigType(storeConf)
	private = false
	return NewCloudStoreWithConfig(config, storeType, private)
}

func GetConfigType(config Config) (cfg interface{}) {
	switch config.CloudType {
	case Cos:
		cfg = &ConfigCos{}
	case Bos:
		cfg = &ConfigBos{}
	case Oss:
		cfg = &ConfigOss{}
	case Minio:
		cfg = &ConfigMinio{}
	case Upyun:
		cfg = &ConfigUpYun{}
	case Qiniu:
		cfg = &ConfigQiniu{}
	case Obs:
		cfg = &ConfigObs{}
	}
	vc := reflect.ValueOf(config)
	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(cfg)
	numFields := t.Elem().NumField()
	for i := 0; i < numFields; i++ {
		fieldName := t.Elem().Field(i).Tag.Get("fieldName")
		val := vc.FieldByName(fieldName)
		v.Elem().FieldByName(fieldName).Set(val)
	}
	return
}

func NewCloudStoreWithConfig(storeConfig interface{}, storeType StoreType, private bool) (cs *CloudStore, err error) {
	var errWithoutConfig = errors.New("云存储配置不正确")
	cs = &CloudStore{
		StoreType: storeType,
		Config:    storeConfig,
	}
	cs.Private = private
	switch cs.StoreType {
	case Oss:
		cfg := cs.Config.(*ConfigOss)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Endpoint == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloud.NewOSS(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint, bucket, domain)
		cs.CanGZIP = true
	case Obs:
		cfg := cs.Config.(*ConfigObs)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Endpoint == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloud.NewOBS(cfg.AccessKey, cfg.SecretKey, bucket, cfg.Endpoint, domain)
	case Qiniu:
		cfg := cs.Config.(*ConfigQiniu)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloud.NewQINIU(cfg.AccessKey, cfg.SecretKey, bucket, domain)
	case Upyun:
		cfg := cs.Config.(*ConfigUpYun)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		//if cfg.Operator == "" || cfg.Password == "" || bucket == "" {
		if cfg.AccessKey == "" || cfg.SecretKey == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		//cs.client = CloudStore2.NewUpYun(bucket, cfg.Operator, cfg.Password, domain, cfg.Secret)
		cs.Client = cloud.NewUpYun(bucket, cfg.AccessKey, cfg.SecretKey, domain, cfg.Endpoint)
	case Minio:
		cfg := cs.Config.(*ConfigMinio)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Endpoint == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloud.NewMinIO(cfg.AccessKey, cfg.SecretKey, bucket, cfg.Endpoint, domain)
		cs.CanGZIP = true
	case Bos:
		cfg := cs.Config.(*ConfigBos)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Endpoint == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloud.NewBOS(cfg.AccessKey, cfg.SecretKey, bucket, cfg.Endpoint, domain)
		cs.CanGZIP = true
	case Cos:
		cfg := cs.Config.(*ConfigCos)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.AppId == "" || bucket == "" || cfg.Region == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloud.NewCOS(cfg.AccessKey, cfg.SecretKey, bucket, cfg.AppId, cfg.Region, domain)
		cs.CanGZIP = true
	}
	return
}
