package service

import (
	"encoding/json"
	"github.com/cilidm/toolbox/gconv"
	"github.com/cilidm/toolbox/levelDB"
	"github.com/cilidm/toolbox/logging"
	conf "ssh_backup/app/config"
	"ssh_backup/app/model"
	"ssh_backup/app/util"
	"sync/atomic"
)

var (
	done        = make(chan bool)
	pathHandle  = make(chan string, conf.Conf.FileInfo.MaxChannel)
	maxReadChan = make(chan bool, conf.Conf.FileInfo.MaxChannel)
	fileIndex   = util.InitialID
	//picIndex    = util.InitialID
	//videoIndex  = util.InitialID
)

// 校验文件是否已入库
func checkFileStatus(v string) bool {
	has, _ := levelDB.GetServer().FindByKey(util.GetLdbKey(conf.Conf.Dir.Source, v))
	if string(has) != "" {
		var oldFile model.FileInfo
		json.Unmarshal(has, &oldFile)
		if oldFile.Status == model.Processed {
			logging.Warn(oldFile.FileName, "已存在")
			return true
		}
	}
	return false
}

func AddLdbIndex(v model.FileInfo) {
	v.Status = model.Processed
	levelDB.GetServer().Insert(util.GetLdbKey(conf.Conf.Dir.Source, v.FileSource), v)

	levelDB.GetServer().Insert(util.GetFileIndexKey(gconv.String(fileIndex)), v)
	atomic.AddUint64(&fileIndex, 1)

	// 是否需要记录图片及视频
	//if str.IsContain(util.PicExt, path.Ext(v.FileName)) {
	//	levelDB.GetServer().Insert(util.GetPicIndexKey(gconv.String(picIndex)), v)
	//	atomic.AddUint64(&picIndex, 1)
	//}

	//if str.IsContain(util.VideoExt, path.Ext(v.FileName)) {
	//	levelDB.GetServer().Insert(util.GetVideoIndexKey(gconv.String(videoIndex)), v)
	//	atomic.AddUint64(&videoIndex, 1)
	//}
}
