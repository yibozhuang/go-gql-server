package main

import (
	"github.com/yibozhuang/go-gql-server/cmd/gql-server/config"
	log "github.com/yibozhuang/go-gql-server/internal/logger"

	"github.com/yibozhuang/go-gql-server/internal/orm"
	"github.com/yibozhuang/go-gql-server/pkg/server"
)

func main() {
	var serverconf = config.Server()

	// Create a new ORM instance
	orm, err := orm.Factory(serverconf)
	defer orm.DB.Close()

	if err != nil {
		log.Panic(err)
	}

	server.Run(serverconf, orm)
}
