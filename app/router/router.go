package router

import (
	"github.com/gin-gonic/gin"
	"ssh_backup/app/controller"
)

func RunServer() {
	r := gin.Default()

	r.LoadHTMLGlob("views/*")

	r.GET("/", controller.Index)		// layout首页

	r.GET("/index", controller.ShowAll)	// 首页展示页

	r.GET("/index_json", controller.ShowAllJson)

	r.Run(":8008")
}
