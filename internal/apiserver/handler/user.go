package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/moweilong/milady/pkg/core"
)

func init() {
	Register(func(v1 *gin.RouterGroup, handler *Handler) {
		// 用户相关路由
		rg := v1.Group("/users")
		rg.POST("", handler.CreateUser) // 创建用户。这里要注意：创建用户是不用进行认证和授权的
		rg.Use(handler.mws...)
		// rg.PUT(":userID/change-password", handler.ChangePassword) // 修改用户密码
		rg.PUT(":userID", handler.UpdateUser)    // 更新用户信息
		rg.DELETE(":userID", handler.DeleteUser) // 删除用户
		rg.GET(":userID", handler.GetUser)       // 查询用户详情
		rg.GET("", handler.ListUser)             // 查询用户列表.
	})
}

// CreateUser handles the creation of a new user.
func (h *Handler) CreateUser(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().Create, h.val.ValidateCreateUserRequest)
}

// UpdateUser handles updating an existing user's details.
func (h *Handler) UpdateUser(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().Update, h.val.ValidateUpdateUserRequest)
}

// DeleteUser handles the deletion of one or more users.
func (h *Handler) DeleteUser(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().Delete, h.val.ValidateDeleteUserRequest)
}

// GetUser retrieves information about a specific user.
func (h *Handler) GetUser(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().Get, h.val.ValidateGetUserRequest)
}

// ListUser retrieves a list of users based on query parameters.
func (h *Handler) ListUser(c *gin.Context) {
	core.HandleJSONRequest(c, h.biz.UserV1().List, h.val.ValidateListUserRequest)
}

// UpdatePassword receives an UpdatePasswordRequest and updates the user's password in the datastore.
// func (h *Handler) UpdatePassword(c *gin.Context) {
// 	core.HandleJSONRequest(c, h.biz.UserV1().UpdatePassword, h.val.ValidateUpdatePasswordRequest)
// }
