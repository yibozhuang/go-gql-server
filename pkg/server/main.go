package server

import (
	"github.com/yibozhuang/go-gql-server/internal/logger"

	"github.com/yibozhuang/go-gql-server/internal/orm"

	"github.com/gin-gonic/gin"
	"github.com/yibozhuang/go-gql-server/pkg/utils"
)

// Run the web server
func Run(serverconf *utils.ServerConfig, orm *orm.ORM) {
	r := gin.Default()

	// Initialize the Auth providers
	InitalizeAuthProviders(serverconf)

	// Routes and Handlers
	RegisterRoutes(serverconf, r, orm)

	// Inform the user where the server is listening
	logger.Info("Running @ " + serverconf.SchemaVersionedEndpoint(""))

	// Run the server
	// Print out and exit(1) to the OS if the server cannot run
	logger.Fatal(r.Run(serverconf.ListenEndpoint()))
}
