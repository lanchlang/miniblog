package util

import "os"

//检查文件或目录是否存在
//输入为：文件或目录的路径
//返回：bool(表示是否存在)，error(调用过程中出现错误)
//如果error不为空，表明有错误
//当error为空时，如果bool为true，表示存在，如果bool为false，表示不存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
