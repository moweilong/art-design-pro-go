package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	v1 "github.com/moweilong/art-design-pro-go/pkg/api/apiserver/v1"
	"github.com/moweilong/milady/pkg/core"
)

func init() {
	Register(func(v1 *gin.RouterGroup, handler *Handler) {
		// 用户相关路由
		rg := v1.Group("/auth")
		rg.POST("/login", handler.Login) // 登录。这里要注意：登录是不用进行认证和授权的
		rg.Use(handler.mws...)
	})
}

// Login authenticates the user credentials and returns a token on success.
func (h *Handler) Login(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.AuthV1().Login, h.val.ValidateLoginRequest)
}

// Logout invalidates the user token.
// func (h *Handler) Logout(c *gin.Context) {
// 	core.HandleJSONRequest(c, h.biz.AuthV1().Logout, h.val.ValidateLogoutRequest)
// }

// RefreshToken generates a new token using the refresh token.
func (h *Handler) RefreshToken(ctx context.Context, rq *v1.RefreshTokenRequest) (*v1.LoginReply, error) {
	return h.biz.AuthV1().RefreshToken(ctx, rq)
}

// Authenticate validates the user token and returns the user ID.
func (h *Handler) Authenticate(ctx context.Context, rq *v1.AuthenticateRequest) (*v1.AuthenticateResponse, error) {
	return h.biz.AuthV1().Authenticate(ctx, rq.Token)
}

// Authorize checks whether the user is authorized for the object/action.
func (h *Handler) Authorize(ctx context.Context, rq *v1.AuthorizeRequest) (*v1.AuthorizeResponse, error) {
	return h.biz.AuthV1().Authorize(ctx, rq.Sub, rq.Obj, rq.Act)
}

// Auth authenticates and authorizes the user token for an object/action.
func (h *Handler) Auth(ctx context.Context, rq *v1.AuthRequest) (*v1.AuthResponse, error) {
	authn, err := h.Authenticate(ctx, &v1.AuthenticateRequest{Token: rq.Token})
	if err != nil {
		return nil, err
	}

	authz, err := h.Authorize(ctx, &v1.AuthorizeRequest{Sub: authn.UserID, Obj: rq.Obj, Act: rq.Act})
	if err != nil {
		return nil, err
	}

	return &v1.AuthResponse{UserID: authn.UserID, Allowed: authz.Allowed}, nil
}
