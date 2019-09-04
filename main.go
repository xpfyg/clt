package main

import (
	"fmt"
	"os"
)

func main() {
	cmdServerV1.Run = serverControllerV1
	cmdServerV2.Run = serverControllerV2
	cmdRunner := &Commands{
		cmdServerV1,
		cmdServerV2,
	}
	cmdRunner.Run()
}

var cmdServerV1 = &Command{
	UsageLine: "cmd1 [jackma | ponyma | start]",
	Short:     "example server",
	Long: `
	jackma              i am jackma.
	ponyma              i am ponyma.
	start      			start the server.
	`,
}

var cmdServerV2 = &Command{
	UsageLine: "server [start | stop | start-console]",
	Short:     "server http",
	Long: `
	start              start the server.
	stop               stop the server.
	start-console      start the server in console.
	`,
}

func serverControllerV1(cmd *Command, args []string) {
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

func serverControllerV2(cmd *Command, args []string) {
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
