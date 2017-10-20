package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//cp命令的用法
const CommandCpUsage = `
  cp 源文件 目标文件
  cp 源文件 目标文件夹/

  eg:
     cp s3://资料库/文件路径 本地文件路径
`

func init() {
	RegisterCommand("cp", CommandCpUsage, CommandCp)
}

//cp命令
func CommandCp(args ...string) {
	if len(args) != 2 {
	}

	src := args[0]
	dst := args[1]

	if strings.HasSuffix(dst, "/") {
		dst += filepath.Base(src)
	}

	fmt.Printf("%s -> %s\n", src, dst)

	fromSeafile := strings.HasPrefix(src, "sf://")
	toSeafile := strings.HasPrefix(dst, "sf://")

	if fromSeafile {
		if toSeafile {
			transfer(src, dst)
		} else {
			download(src, dst)
		}
	} else {
		if toSeafile {
			upload(src, dst)
		} else {
			fmt.Println("不支持本地复制")
		}
	}
}

//上传到Seafile资料库
func upload(src, dst string) {
	libName, dst := parseDirectory(strings.TrimPrefix(dst, "sf://"))
	library, err := sf.GetLibrary(libName)
	if err != nil {
		fmt.Println("获取资料库失败: %s", err)
		return
	}

	b, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println("无法读取源文件: %s", err)
		return
	}

	dir := filepath.Dir(dst) + "/"
	fname := filepath.Base(dst)
	contents := map[string][]byte{fname: b}
	err = library.UploadFileContent(dir, contents)
	if err != nil {
		fmt.Println("上传失败: %s", err)
	}
}

//从Seafile资料库下载文件
func download(src, dst string) {
	libName, src := parseDirectory(strings.TrimPrefix(src, "sf://"))
	library, err := sf.GetLibrary(libName)
	if err != nil {
		fmt.Println("获取资料库失败: %s", err)
		return
	}

	b, err := library.FetchFileContent(src)
	if err != nil {
		fmt.Println("下载失败: %s", err)
		return
	}

	file, err := os.Create(dst)
	if err != nil {
		fmt.Println("打开文件失败: %s", err)
		return
	}
	defer file.Close()

	_, err = file.Write(b)
	if err != nil {
		fmt.Println("写入文件失败: %s", err)
		return
	}
}

//Seafile资料库之间复制文件
func transfer(src, dst string) {
	srcLibName, src := parseDirectory(strings.TrimPrefix(src, "sf://"))
	srcLibrary, err := sf.GetLibrary(srcLibName)
	if err != nil {
		fmt.Println("获取源资料库失败: %s", err)
		return
	}

	fmt.Println(srcLibrary)

	dstLibName, dst := parseDirectory(strings.TrimPrefix(dst, "sf://"))
	dstLibrary, err := sf.GetLibrary(dstLibName)
	if err != nil {
		fmt.Println("获取目标资料库失败: %s", err)
		return
	}

	fmt.Println(dstLibrary)

	err = srcLibrary.CopyFileToLibrary(src, dstLibrary.Id, filepath.Dir(dst))
	if err != nil {
		fmt.Println("复制文件失败: %s", err)
		return
	}
}
