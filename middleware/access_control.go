package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"go-web-demo/component"
	"log"
	"net/http"
)

// DefaultAuthorize determines if current subject has been authorized to take an action on an object.
func DefaultAuthorize(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get current user/subject
		token := c.Request.Header.Get("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, component.RestResponse{Message: "token is nil"})
			return
		}
		username, err := component.GlobalCache.Get(token)
		if err != nil || string(username) == "" {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, component.RestResponse{Message: "user hasn't logged in yet"})
			return
		}

		// Casbin enforces policy
		ok, err := enforce(string(username), obj, act, component.Enforcer)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, component.RestResponse{Message: "error occurred when authorizing user"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, component.RestResponse{Message: "forbidden"})
			return
		}

		c.Next()
	}
}

func enforce(sub string, obj string, act string, enforcer *casbin.Enforcer) (bool, error) {
	// Load policies from DB dynamically
	err := enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	// Verify
	ok, err := enforcer.Enforce(sub, obj, act)
	return ok, err
}

func AuthorizeAdapterAndModel(obj string, act string, adapter *gormadapter.Adapter, model string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get current user/subject
		token := c.Request.Header.Get("token")
		if token == "" {
			c.AbortWithStatusJSON(401, component.RestResponse{Message: "token is nil"})
			return
		}
		username, err := component.GlobalCache.Get(token)
		if err != nil || string(username) == "" {
			log.Println(err)
			c.AbortWithStatusJSON(401, component.RestResponse{Message: "user hasn't logged in yet"})
			return
		}

		// Load model configuration file and policy store adapter
		enforcer, err := casbin.NewEnforcer(model, adapter)
		// Casbin enforces policy
		ok, err := enforce(string(username), obj, act, enforcer)

		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, component.RestResponse{Message: "error occurred when authorizing user"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, component.RestResponse{Message: "forbidden"})
			return
		}

		c.Next()
	}
}
