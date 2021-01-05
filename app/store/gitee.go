package store

/**
官方文档地址 https://gitee.com/api/v5/swagger#/deleteV5ReposOwnerRepoContentsPath
*/

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cilidm/toolbox/logging"
	"github.com/kirinlabs/HttpRequest"
	"io/ioutil"
	"os"
	"ssh_backup/app/config"
)

var (
	req      = HttpRequest.NewRequest()
	token    = config.Conf.Gitee.Token
	owner    = config.Conf.Gitee.Owner
	repo     = config.Conf.Gitee.Repo
	apiUrl   = "https://gitee.com/api/v5/repos/"
	contents = apiUrl + owner + "/" + repo + "/contents/"
	blob     = apiUrl + owner + "/" + repo + "/git/blobs/"
)

type Contents struct {
	Links       link   `json:"_links"`
	Content     string `json:"content"`
	DownloadUrl string `json:"download_url"`
	Encoding    string `json:"encoding"`
	HtmlUrl     string `json:"html_url"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        string `json:"size"`
	Type        string `json:"type"`
	Url         string `json:"url"`
}

type link struct {
	Self string `json:"self"`
	Html string `json:"html"`
}

// 获取仓库具体路径下的内容
func GetContents(dir string) (cons []Contents, err error) {
	resp, err := req.Get(contents + dir)
	if err != nil {
		logging.Error(err.Error())
		return cons, err
	}
	defer resp.Close()
	fmt.Println(resp.StatusCode())
	body, err := resp.Body()
	if err != nil {
		logging.Error(err.Error())
		return cons, err
	}
	fmt.Println(string(body))
	err = json.Unmarshal(body, &cons)
	return cons, err
}

// 通用get
func HttpGet(url string, data map[string]interface{}) (body []byte, err error) {
	resp, err := req.Get(url, data)
	if err != nil {
		logging.Error(err.Error())
		return body, err
	}
	defer resp.Close()
	body, err = resp.Body()
	if err != nil {
		logging.Error(err.Error())
	}
	return body, err
}

// 上传文件
// fileName : /file/main.go
func CreateFile(filePath, fileName, msg string) error {
	req.SetHeaders(map[string]string{
		"Content-Type": "application/json;charset=UTF-8",
		"Connection":   "keep-alive",
	})
	file := FileToBase64(filePath)
	resp, err := req.Post(contents+fileName, map[string]interface{}{
		"access_token": token,
		"content":      file,
		"message":      msg,
	})
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	defer resp.Close()
	body, err := resp.Body()
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	fmt.Println(string(body), err) // 可以获取sha
	return nil
}

// sha
func DelFile(fileName, sha, msg string) error {
	resp, err := req.Delete(contents+fileName, map[string]interface{}{
		"access_token": token,
		"sha":          sha,
		"message":      msg,
	})
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	defer resp.Close()
	body, err := resp.Body()
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	fmt.Println(string(body), err)
	return nil
}

// 更新文件
func UpdateFile(filePath, fileName, sha, msg string) error {
	file := FileToBase64(filePath)
	resp, err := req.Put(contents+fileName, map[string]interface{}{
		"access_token": token,
		"content":      file,
		"sha":          sha,
		"message":      msg,
	})
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	defer resp.Close()
	body, err := resp.Body()
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	fmt.Println(string(body), err)
	return nil
}

type FileBlob struct {
	Sha      string `json:"sha"`
	Size     int    `json:"size"`
	Encoding string `json:"encoding"`
	Url      string `json:"url"`
	Content  string `json:"content"`
}

func GetBlob(sha string) (FileBlob, error) {
	b, err := HttpGet(blob+sha, nil)
	if err != nil {
		logging.Error(err)
	}
	var fb FileBlob
	json.Unmarshal(b, &fb)
	return fb, nil
}

func FileToBase64(filePath string) string {
	ff, _ := os.Open(filePath)
	defer ff.Close()
	sourcebuffer := make([]byte, 500000)
	n, _ := ff.Read(sourcebuffer)
	sourcestring := base64.StdEncoding.EncodeToString(sourcebuffer[:n])
	return sourcestring
}

func Base64ToFile(encode, saveFile string) error {
	b, err := base64.StdEncoding.DecodeString(encode) //成图片文件并把文件写入到buffer
	if err != nil {
		logging.Error(err)
		return err
	}
	err = ioutil.WriteFile(saveFile, b, os.ModePerm) //buffer输出到文件中（不做处理，直接写到文件）
	return err
}
