package service

import (
	"os"
	"path"
	"path/filepath"

	"github.com/cilidm/toolbox/logging"
	"github.com/pkg/sftp"
)

// remotePath服务器路径 localPath本地路径
func GetFile(sftpClient *sftp.Client, remotePath, localPath string) error {
	srcFile, err := sftpClient.Open(remotePath)
	if err != nil {
		logging.Fatal(err)
	}
	defer srcFile.Close()
	dstFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		return err
	}
	logging.Info("copy ", srcFile.Name(), " finished!")
	return nil
}

// remoteDir目标文件文件夹名称,remoteName目标文件名, localFilePath源文件夹路径
// sourceClient源机器连接  targetClient目标机器连接
func UploadFile(sourceClient *sftp.Client, targetClient *sftp.Client, remoteDir, remoteName, localFilePath string, localSize int64) error {
	has, err := targetClient.Stat(filepath.Join(remoteDir, remoteName))
	if err == nil && (has.Size() == localSize) {
		logging.Warn("文件", remoteName, "已s存在")
		return nil
	}
	_, err = ClientPathExists(remoteDir, targetClient)
	if err != nil {
		return err
	}
	srcFile, err := sourceClient.Open(localFilePath)
	if err != nil {
		logging.Error("源文件无法读取", err.Error())
		return err
	}
	defer srcFile.Close()
	dstFile, err := targetClient.Create(path.Join(remoteDir, remoteName)) // 如果文件存在，create会清空原文件 openfile会追加
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 10000)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[:n]) // 读多少 写多少
	}
	logging.Info(localFilePath, "传输成功")
	return nil
}

func UploadFromLocal(targetClient *sftp.Client, remoteDir, remoteName, localFilePath string, localSize int64) error {
	has, err := targetClient.Stat(filepath.Join(remoteDir, remoteName))
	if err == nil && (has.Size() == localSize) {
		logging.Warn("文件", remoteName, "已s存在")
		return nil
	}
	targetClient.MkdirAll(remoteDir)
	err = targetClient.Chmod(remoteDir, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = ClientPathExists(remoteDir, targetClient)
	if err != nil {
		return err
	}

	srcFile, err := os.Open(localFilePath)
	if err != nil {
		logging.Error("源文件无法读取", err.Error())
		return err
	}
	defer srcFile.Close()
	dstFile, err := targetClient.Create(filepath.Join(remoteDir, remoteName)) // 如果文件存在，create会清空原文件 openfile会追加
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 10000)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[:n]) // 读多少 写多少
	}
	logging.Info(localFilePath, "传输成功")
	return nil
}

func ClientPathExists(path string, client *sftp.Client) (bool, error) {
	_, err := client.Stat(path)
	if err == nil {
		return true, nil
	}
	err = client.MkdirAll(path)
	if err != nil {
		return false, err
	}
	err = client.Chmod(path, os.ModePerm)
	return false, err
}
