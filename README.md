# clt
### 描述
基于go语言的命令行工具，根据输入提示可执行程序的帮助信息，适合用于复杂业务的脚本辅助

### 使用效果
```
$ ./test help
commands server - connect everysites!

Usage:

    test command [arguments]

The commands are:

    cmd1        example server
    server      server http

Use "test help [command]" for more information about a command.

$ ./test help cmd1
usage: test cmd1 [jackma | ponyma | start]

        jackma              i am jackma.
        ponyma              i am ponyma.
        start               start the server.

$ ./test cmd1 jackma
悔创阿里杰克马
```

### 如何使用
```
package main

import (
	"fmt"
	"os"

	"github.com/xpfyg/clt"
)

func main() {
	cmdServerV1.Run = serverControllerV1
	cmdServerV2.Run = serverControllerV2
	cmdRunner := clt.Commands{
		CommandList: []*clt.Command{
			cmdServerV1,
			cmdServerV2,
		},
		ApiName: "test",
	}

	cmdRunner.Run()
}

var cmdServerV1 = &clt.Command{
	UsageLine: "cmd1 [jackma | ponyma | start]",
	Short:     "example server",
	Long: `
	jackma              i am jackma.
	ponyma              i am ponyma.
	start               start the server.
	`,
}

var cmdServerV2 = &clt.Command{
	UsageLine: "server [start | stop | start-console]",
	Short:     "server http",
	Long: `
	start              start the server.
	stop               stop the server.
	start-console      start the server in console.
	`,
}

func serverControllerV1(cmd *clt.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("what's your problem!")
		os.Exit(2)
	}

	switch args[0] {
	case "jackma":
		fmt.Println("悔创阿里杰克马")
	case "ponyma":
		fmt.Println("普通家庭马化腾")
	case "start":
		fmt.Println("开始你的表演")
	default:
		fmt.Println("what's your problem!")
	}
}

func serverControllerV2(cmd *clt.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("what's your problem!")
		os.Exit(2)
	}

	switch args[0] {
	case "start":
		fmt.Println("server start")
	case "stop":
		fmt.Println("server stop")
	case "start-console":
		fmt.Println("server start in console")
	default:
		fmt.Println("what's your problem!")
	}
}
```
