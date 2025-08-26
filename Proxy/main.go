package main

import (
	"github.com/akmalfairuz/legacy-version/legacyver"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"proxy/proxy"
)

func main() {
	config := proxy.ReadConfig()
	token, err := auth.RequestLiveToken()
	if err != nil {
		panic(err)
	}
	src := auth.RefreshTokenSource(token)

	p, err := minecraft.NewForeignStatusProvider(config.Connection.RemoteAddress)
	if err != nil {
		panic(err)
	}
	listener, err := minecraft.ListenConfig{
		StatusProvider:       p,
		AcceptedProtocols:    legacyver.All(false),
		TexturePacksRequired: true,
		ResourcePacks:        proxy.ParsePacks(),
	}.Listen("raknet", config.Connection.LocalAddress)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		c, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go proxy.HandleConn(c.(*minecraft.Conn), listener, config, src)
	}
}
