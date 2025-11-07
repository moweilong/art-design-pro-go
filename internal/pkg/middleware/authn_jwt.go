package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"

	"github.com/moweilong/art-design-pro-go/internal/apiserver/model"
	"github.com/moweilong/art-design-pro-go/internal/pkg/contextx"
	"github.com/moweilong/art-design-pro-go/internal/pkg/errno"
	"github.com/moweilong/milady/pkg/authn"
	"github.com/moweilong/milady/pkg/core"
)

const (
	// reason 错误原因
	reason string = "UNAUTHORIZED"

	// bearerWord Bearer关键字
	bearerWord string = "Bearer"

	// authorizationKey 请求头中的Authorization键
	authorizationKey string = "Authorization"
)

var (
	// ErrMissingJwtToken JWT令牌缺失错误
	ErrMissingJwtToken = errors.Unauthorized(reason, "JWT token is missing")
)

// UserRetriever 用于根据用户名获取用户的接口.
type UserRetriever interface {
	// GetUser 根据用户ID获取用户信息
	GetUser(ctx context.Context, userID string) (*model.UserM, error)
}

// AuthnJWT 是Gin框架的JWT认证中间件
// 功能与Kratos版本的Server中间件相同，但适配Gin框架
func AuthnMiddleware(a authn.Authenticator, retriever UserRetriever) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader(authorizationKey)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "JWT token is missing",
			})
			return
		}

		// 解析Bearer Token
		auths := strings.SplitN(authHeader, " ", 2)
		if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid authorization format",
			})
			return
		}

		accessToken := auths[1]
		ctx := c.Request.Context()

		// 解析JWT声明
		claims, err := a.ParseClaims(ctx, accessToken)
		if err != nil {
			// 处理错误情况
			statusCode := http.StatusUnauthorized
			message := "Invalid token"

			// 根据错误类型提供更具体的错误信息
			if se, ok := err.(*errors.Error); ok {
				message = se.Message
			}

			c.AbortWithStatusJSON(statusCode, gin.H{
				"code":    statusCode,
				"message": message,
			})
			return
		}

		user, err := retriever.GetUser(c, claims.Subject)
		if err != nil {
			core.WriteResponse(c, nil, errno.ErrUnauthenticated.WithMessage("%s", err.Error()))
			c.Abort()
			return
		}

		// 将信息注入上下文
		// 1. 注入到Gin上下文
		c.Set("userID", user.UserID)
		c.Set("accessToken", accessToken)
		c.Set("claims", claims)

		// 2. 同时更新请求上下文，保持与原有逻辑一致
		newCtx := contextx.WithClaims(ctx, claims)
		newCtx = contextx.WithUserID(newCtx, user.UserID)
		newCtx = contextx.WithAccessToken(newCtx, accessToken)
		c.Request = c.Request.WithContext(newCtx)

		c.Next()
	}
}

// AuthnJWTSkip 允许跳过某些路径的JWT认证
func AuthnJWTSkip(a authn.Authenticator, retriever UserRetriever, skipPaths ...string) gin.HandlerFunc {
	// 创建跳过路径的映射
	skipPathMap := make(map[string]struct{})
	for _, path := range skipPaths {
		skipPathMap[path] = struct{}{}
	}

	return func(c *gin.Context) {
		// 检查当前路径是否需要跳过认证
		if _, skip := skipPathMap[c.Request.URL.Path]; skip {
			c.Next()
			return
		}

		// 否则执行认证中间件
		authnJWT := AuthnMiddleware(a, retriever)
		authnJWT(c)
	}
}

// GetUserID 从Gin上下文获取用户ID的辅助函数
func GetUserID(c *gin.Context) string {
	userID, exists := c.Get("userID")
	if !exists {
		return ""
	}
	return userID.(string)
}

// GetAccessToken 从Gin上下文获取访问令牌的辅助函数
func GetAccessToken(c *gin.Context) string {
	token, exists := c.Get("accessToken")
	if !exists {
		return ""
	}
	return token.(string)
}
