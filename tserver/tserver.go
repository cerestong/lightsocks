package tserver

import (
	"github.com/cerestong/lightsocks/core"
	"log"
	"net"
)

// LsTServer 远端代理实体
type LsTServer struct {
	*core.SecureSocket
}

// New 新建一个远端代理
// 监听外部请求
// 请求转发前解密数据
// 应答转发前加密数据
func New(password *core.Password, listenAddr, remoteAddr *net.TCPAddr) *LsTServer {
	return &LsTServer{
		SecureSocket: &core.SecureSocket{
			Cipher: core.NewCipher(password),
			ListenAddr: listenAddr,
			RemoteAddr: remoteAddr,
		},
	}
}

// Listen 启动监听
func (local *LsTServer) Listen(didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", local.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	if didListen != nil {
		didListen(listener.Addr())
	}

	for {
		userConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		// userConn被关闭时直接清除所有数据 丢弃没有发送的数据
		userConn.SetLinger(0)
		go local.handleConn(userConn)
	}
	
	return nil
}

func (local *LsTServer) handleConn(userConn *net.TCPConn) {
	defer userConn.Close()

	proxyServer, err := local.DialRemote()
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyServer.Close()
	proxyServer.SetLinger(0)

	// 应答加密转发
	go func() {
		err := local.EncodeCopy(userConn, proxyServer)
		if err != nil {
			userConn.Close()
			proxyServer.Close()
		}
	}()

	// 进行解密转发
	local.DecodeCopy(proxyServer, userConn)
}