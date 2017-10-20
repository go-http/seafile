package main

//命令具体执行的函数
type CommandFunc func(...string)

//用于注册的命令结构
type Command struct {
	Usage string
	Func  CommandFunc
}

var commandMap = map[string]Command{}

//注册新命令
func RegisterCommand(cmd, usage string, f CommandFunc) {
	commandMap[cmd] = Command{
		Usage: usage,
		Func:  f,
	}
}
