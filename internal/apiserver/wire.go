//go:build wireinject
// +build wireinject

package apiserver

import (
	"github.com/google/wire"
	"github.com/moweilong/milady/pkg/authz"
	genericoptions "github.com/moweilong/milady/pkg/options"

	"github.com/moweilong/art-design-pro-go/internal/apiserver/biz"
	"github.com/moweilong/art-design-pro-go/internal/apiserver/pkg/validation"
	"github.com/moweilong/art-design-pro-go/internal/apiserver/store"
	"github.com/moweilong/art-design-pro-go/internal/pkg/auth"
	mw "github.com/moweilong/art-design-pro-go/internal/pkg/middleware/gin"
)

// NewServer sets up and create the web server with all necessary dependencies.
func NewServer(*Config, *genericoptions.JWTOptions, *genericoptions.RedisOptions) (*Server, error) {
	wire.Build(
		NewWebServer,
		wire.Struct(new(ServerConfig), "*"), // * 表示注入全部字段
		wire.Struct(new(Server), "*"),
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		store.SetterProviderSet, // auth 依赖的 store.Setter
		ProvideDB,               // 提供数据库实例
		NewAuthenticator,        // 提供认证器
		auth.ProviderSet,
		validation.ProviderSet,
		wire.NewSet(
			wire.Struct(new(UserRetriever), "*"),
			wire.Bind(new(mw.UserRetriever), new(*UserRetriever)),
		),
		authz.ProviderSet,
	)
	return nil, nil
}
