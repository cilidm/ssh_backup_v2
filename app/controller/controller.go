package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ssh_backup/app/controller/api"
	"ssh_backup/app/model"
	"ssh_backup/app/service"
	"strconv"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func ShowAll(c *gin.Context) { // 全部文件
	c.HTML(http.StatusOK, "main.html", gin.H{})
}

func ShowAllJson(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	resp, count := service.GetAllFileByPage(page, limit)
	var data []model.FileInfo
	for _, v := range resp {
		data = append(data, v.(model.FileInfo))
	}
	api.SuccessResp(c).SetCode(0).SetData(data).SetCount(count).WriteJsonExit()
}