package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yibozhuang/go-gql-server/internal/handlers"
	auth "github.com/yibozhuang/go-gql-server/internal/handlers/auth/middleware"
	"github.com/yibozhuang/go-gql-server/internal/logger"
	"github.com/yibozhuang/go-gql-server/internal/orm"
	"github.com/yibozhuang/go-gql-server/pkg/utils"
)

// GraphQL routes
func GraphQL(cfg *utils.ServerConfig, r *gin.Engine, orm *orm.ORM) error {
	// GraphQL paths
	gqlPath := cfg.VersionedEndpoint(cfg.GraphQL.Path)
	pgqlPath := cfg.GraphQL.PlaygroundPath
	g := r.Group(gqlPath)

	// GraphQL handler
	g.POST("", auth.Middleware(g.BasePath(), cfg, orm), handlers.GraphqlHandler(orm, &cfg.GraphQL))
	logger.Info("GraphQL @ ", gqlPath)
	// Playground handler
	if cfg.GraphQL.IsPlaygroundEnabled {
		logger.Info("GraphQL Playground @ ", g.BasePath()+pgqlPath)
		g.GET(pgqlPath, handlers.PlaygroundHandler(g.BasePath()))
	}

	return nil
}
