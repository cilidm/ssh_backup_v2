package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommonResp struct {
	Code  int           `json:"code"` //响应编码 0 成功 500 错误 403 无权限  -1  失败
	Msg   string        `json:"msg"`  //消息
	Data  interface{}   `json:"data"` //数据内容
	Count int           `json:"count,omitempty"`
}

type Resp struct {
	c *gin.Context
	r *CommonResp
}

//返回一个成功的消息体
func SuccessResp(c *gin.Context) *Resp {
	msg := CommonResp{
		Code: 200,
		Msg:  "操作成功",
	}
	var a = Resp{
		r: &msg,
		c: c,
	}
	return &a
}

//返回一个错误的消息体
func ErrorResp(c *gin.Context) *Resp {
	msg := CommonResp{
		Code: 500,
		Msg:  "操作失败",
	}
	var a = Resp{
		r: &msg,
		c: c,
	}
	return &a
}

//设置消息体的内容
func (resp *Resp) SetMsg(msg string) *Resp {
	resp.r.Msg = msg
	return resp
}

//设置消息体的编码
func (resp *Resp) SetCode(code int) *Resp {
	resp.r.Code = code
	return resp
}

//设置消息体的数据
func (resp *Resp) SetData(data interface{}) *Resp {
	resp.r.Data = data
	return resp
}

//设置消息体的业务类型
func (resp *Resp) SetCount(count int) *Resp {
	resp.r.Count = count
	return resp
}

//输出json到客户端
func (resp *Resp) WriteJsonExit() {
	resp.c.JSON(http.StatusOK, resp.r)
	resp.c.Abort()
}

func (resp *Resp) WriteErrJsonExit(errCode int) {
	resp.c.JSON(errCode, resp.r)
	resp.c.Abort()
}
