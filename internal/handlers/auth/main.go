package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yibozhuang/go-gql-server/pkg/utils"
)

func addProviderToContext(c *gin.Context, value interface{}) *http.Request {
	return c.Request.WithContext(context.WithValue(c.Request.Context(),
		string(utils.ProjectContextKeys.ProviderCtxKey), value))
}
