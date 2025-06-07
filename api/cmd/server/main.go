package main

import (
	"github.com/newssourcecrawler/realtorinstall/api/internal/server"
)

func main() {
	server.Start() // handles openDB(), migrations, repos, services, routes, TLS & graceful shutdown
}
