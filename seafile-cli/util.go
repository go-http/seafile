package main

import (
	"fmt"
	"strings"
)

//解析文件夹参数
//    以/开头的参数，表示默认资料库的下文件夹完整路径
//    非/开头的参数，/前表示资料库名，/及之后表示文件夹的完整路径
func parseDirectory(directory string) (string, string) {
	if strings.HasPrefix(directory, "/") {
		return "", directory
	} else {
		strs := strings.SplitN(directory, "/", 2)
		return strs[0], "/" + strs[1]
	}
}

const (
	KiB = 1024
	MiB = 1024 * KiB
	GiB = 1024 * MiB
	TiB = 1024 * GiB
	PiB = 1024 * TiB
	EiB = 1024 * PiB
)

//将int类型的size转换为K、M、B这类格式的字符串
func humanSize(size int) string {
	switch {
	case size > KiB:
		return fmt.Sprintf("%.2fK", float32(size)/KiB)
	case size > MiB:
		return fmt.Sprintf("%.2fM", float32(size)/MiB)
	case size > GiB:
		return fmt.Sprintf("%.2fG", float32(size)/GiB)
	case size > TiB:
		return fmt.Sprintf("%.2fT", float32(size)/TiB)
	case size > PiB:
		return fmt.Sprintf("%.2fP", float32(size)/PiB)
	case size > EiB:
		return fmt.Sprintf("%.2fE", float32(size)/EiB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}
