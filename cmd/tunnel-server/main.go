package main

import (
	"github.com/phayes/freeport"
	"fmt"
	"github.com/cerestong/lightsocks/cmd"
	"github.com/cerestong/lightsocks/core"
	"github.com/cerestong/lightsocks/tserver"
	"log"
	"net"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)

	// 服务器监听端口随机生成
	port, err := freeport.GetFreePort()
	if err != nil {
		port = 7558
	}
	config := &cmd.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		// 密码随机生成
		Password: core.RandPassword().String(),
	}
	config.ReadConfig()
	config.SaveConfig()

	// 解析配置
	password, err := core.ParsePassword(config.Password)
	if err != nil {
		log.Fatalln(err)
	}
	listenAddr, err := net.ResolveTCPAddr("tcp", config.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	remoteAddr, err := net.ResolveTCPAddr("tcp", config.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}
	
	// 启动server端监听
	lsLocal := tserver.New(password, listenAddr, remoteAddr)
	log.Fatalln(lsLocal.Listen(func(listenAddr net.Addr) {
		log.Println("使用配置：", fmt.Sprintf(`
本地监听地址 listen:
%s
远程服务地址 remote:
%s
密码 password:
%s
			`, listenAddr, remoteAddr, password))
		log.Printf("tunnel-server:%s 启动成功 监听在%s\n",
			version, listenAddr.String())
	}))
}