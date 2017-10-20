package main

import (
	"fmt"
	"os"
	"time"
)

//ls命令的用法
const CommandLsUsage = `
  ls               查看资料库列表
  ls /路径         查看默认资料库指定路径的内容
  ls 资料库名/路径 查看指定资料库指定路径的内容
`

func init() {
	RegisterCommand("ls", CommandLsUsage, CommandLs)
}

//ls命令
func CommandLs(args ...string) {
	//不提供文件夹路径，则获取资料库列表
	if len(args) == 0 {
		libraries, err := sf.ListAllLibraries()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, library := range libraries {
			fmt.Println(library.Name)
		}
		return
	}

	//提取资料库名和文件夹路径
	libName, dir := parseDirectory(args[0])

	//获取资料库
	library, err := sf.GetLibrary(libName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取资料库失败%s:\n", err)
		return
	}

	//获取文件夹内容
	entries, err := library.ListDirectoryEntries(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "获取文件夹内容失败%s:\n", err)
		return
	}

	//输出文件夹内容
	fmt.Printf("%s%s 中有%d个项目\n", libName, dir, len(entries))
	for _, e := range entries {
		t := time.Unix(int64(e.Mtime), 0).Format("2006-01-02 15:04:05")
		name := e.Name
		if e.Type == "dir" {
			name += "/"
		}
		fmt.Printf("  %s-%s %s %7s %s\n", e.Type[0:1], e.Permission, t, humanSize(e.Size), name)
	}
}
