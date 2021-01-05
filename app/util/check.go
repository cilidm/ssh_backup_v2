package util

import "os"

func Find(what int, where []int) (idx int) {
	for i, v := range where {
		if v == what {
			return i
		}
	}
	return -1
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}