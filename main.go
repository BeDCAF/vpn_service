package main

import (
	"log"

	"github.com/BeDCAF/vpn_service/config"
	"github.com/BeDCAF/vpn_service/panel"
	"github.com/BeDCAF/vpn_service/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	client, err := panel.NewClientVPN(cfg.Host, cfg.Token)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	httpHandlers := server.NewHTTPHandlers(client)
	httpServer := server.NewHTTPServer(httpHandlers)

	if err := httpServer.StartServer(cfg.Addr, cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
