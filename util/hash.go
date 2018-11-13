package util

import (
	"crypto/sha1"
	"encoding/hex"
)

//使用sha1算法获取字节数组的hash值
//输入为字节数组
//输出为字节数组经过sha1运算后的hash字符串
func HashBytes(bytes []byte) string {
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write(bytes)
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要。
	return hex.EncodeToString(h.Sum(nil))
}
//使用sha1算法获取字符串的hash值
//输入为字符串
//输出为字符串经过sha1运算后的hash字符串
func Hash(str string) string {
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(str))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要。
	return hex.EncodeToString(h.Sum(nil))
}
