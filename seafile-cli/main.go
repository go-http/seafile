//命令行版本的Seafile
package main

import (
	".."
	"flag"
	"fmt"
	"os"
)

var sf *seafile.Client

func main() {
	var host, user, pass, token string
	flag.StringVar(&host, "h", "", "Seafile服务器地址")
	flag.StringVar(&user, "u", "", "Seafile服务器用户名")
	flag.StringVar(&pass, "p", "", "Seafile服务器密码")
	flag.StringVar(&token, "token", "", "Seafile服务器AuthToken")

	flag.Parse()

	if host == "" {
		host = os.Getenv("SEAFILE_HOST")
	}
	if user == "" {
		user = os.Getenv("SEAFILE_USER")
	}
	if pass == "" {
		pass = os.Getenv("SEAFILE_PASS")
	}
	if token == "" {
		token = os.Getenv("SEAFILE_TOKEN")
	}

	//创建客户端
	if token != "" {
		sf = seafile.New(host, token)
	} else {
		sf = seafile.New(host, user, pass)
	}

	cmd, found := commandMap[flag.Arg(0)]
	if !found {
		fmt.Fprintf(os.Stderr, "Usage of %s: %s [选项] <命令> [参数]\n", os.Args[0], os.Args[0])
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		for name, cmd := range commandMap {
			fmt.Fprintf(os.Stderr, "命令%s:%s\n", name, cmd.Usage)
		}
		return
	}

	args := flag.Args()[1:]
	cmd.Func(args...)
}
