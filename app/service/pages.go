package service

import (
	"encoding/json"
	"github.com/cilidm/toolbox/gconv"
	"github.com/cilidm/toolbox/levelDB"
	"github.com/cilidm/toolbox/logging"
	"ssh_backup/app/model"
	"ssh_backup/app/util"
)

func GetAllFileByPage(page, limit int) ([]interface{}, int) {
	count, _ := levelDB.GetServer().FindByPrefix(util.FileIndexKey)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	begin := util.InitialID + uint64((page-1)*limit) // InitialID - 1
	end := begin + uint64(limit)
	if end > util.InitialID+uint64(len(count))-1 {
		end = util.InitialID + uint64(len(count))
	}
	datas, _ := levelDB.GetServer().FindLimit(util.GetFileIndexKey(gconv.String(begin)), util.GetFileIndexKey(gconv.String(end)))
	var resp []interface{}
	for _, v := range datas {
		var f model.FileInfo
		json.Unmarshal([]byte(v), &f)
		resp = append(resp, f)
	}
	return resp, len(count)
}

// 按长度截取字符串 超出的部分用..替代 （发现直接用css样式控制更简单，此方法弃用）
func OmitStr(str string, num int) (newStr string) {
	runeStr := []rune(str)
	if len(runeStr) > num {
		newStr = string(runeStr[:num]) + ".."
	} else {
		newStr = str
	}
	return
}

func GetIndexDetail(key string) (f model.FileInfo) {
	data, err := levelDB.GetServer().FindByKey(key)
	if err != nil {
		logging.Error(err.Error())
		return f
	}
	err = json.Unmarshal(data, &f)
	if err != nil {
		logging.Error(err.Error())
		return f
	}
	return f
}
