package main

import (
	"fmt"
	"imageSync/src"
	"runtime"
)

func main() {

	//命令行参数（含子命令的实例化）
	imageName := src.NewCmdParams()
	if len(imageName) == 0 {
		src.Exit(0, "init命令执行完成")
	}

	//判断当前系统类型
	switch runtime.GOOS {
	case "linux":
		//TODO: 这里没有判断Linux上的docker服务是否运行，后期需要加上
		fmt.Println("当前操作系统为：Linux")
		src.Pull(imageName)
		src.Push(imageName)

	case "windows":
		//TODO: 这里没有判断Windows上的docker服务是否运行，后期需要加上
		fmt.Println("当前操作系统为：Windows")
		src.Pull(imageName)
		src.Push(imageName)

	case "darwin":
		fmt.Println("当前操作系统为：MacOS")
		if src.MacDockerRunCheck() { //判断Docker daemon是否启动
			src.Pull(imageName) //拉取镜像
			src.Push(imageName) //推送镜像
		} else {
			src.Exit(2, "docker daemon未启动")
		}

	default:
		src.Exit(1, "无法判断当前操作系统类型，imageSync停止运行")
	}
}
