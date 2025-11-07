package apiserver

import (
	"context"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/moweilong/milady/pkg/authn"
	jwtauthn "github.com/moweilong/milady/pkg/authn/jwt"
	jwtredis "github.com/moweilong/milady/pkg/authn/jwt/store/redis"
	"github.com/moweilong/milady/pkg/authz"
	genericoptions "github.com/moweilong/milady/pkg/options"
	"github.com/moweilong/milady/pkg/server"
	"github.com/moweilong/milady/pkg/store/registry"
	"github.com/moweilong/milady/pkg/store/where"
	"gorm.io/gorm"

	"github.com/moweilong/art-design-pro-go/internal/apiserver/biz"
	"github.com/moweilong/art-design-pro-go/internal/apiserver/model"
	"github.com/moweilong/art-design-pro-go/internal/apiserver/pkg/validation"
	"github.com/moweilong/art-design-pro-go/internal/apiserver/store"
	"github.com/moweilong/art-design-pro-go/internal/pkg/contextx"
	mw "github.com/moweilong/art-design-pro-go/internal/pkg/middleware"
)

// Config contains application-related configurations.
type Config struct {
	TLSOptions   *genericoptions.TLSOptions
	HTTPOptions  *genericoptions.HTTPOptions
	MySQLOptions *genericoptions.MySQLOptions
	JWTOptions   *genericoptions.JWTOptions
	RedisOptions *genericoptions.RedisOptions
}

// Server represents the web server.
type Server struct {
	cfg *ServerConfig
	srv server.Server
}

// ServerConfig contains the core dependencies and configurations of the server.
type ServerConfig struct {
	*Config
	biz       biz.IBiz
	val       *validation.Validator
	retriever mw.UserRetriever
	authz     *authz.Authz
}

// NewServer initializes and returns a new Server instance.
func (cfg *Config) NewServer(ctx context.Context) (*Server, error) {
	where.RegisterTenant("userID", func(ctx context.Context) string {
		return contextx.UserID(ctx)
	})

	// 初始化 token 包的签名密钥、认证 Key 及 Token 默认过期时间
	// token.Init(cfg.JWTKey, token.WithIdentityKey(known.XUserID), token.WithExpiration(cfg.Expiration))
	// Create the core server instance.
	return NewServer(cfg, cfg.JWTOptions, cfg.RedisOptions)
}

// Run starts the server and listens for termination signals.
// It gracefully shuts down the server upon receiving a termination signal.
func (s *Server) Run(ctx context.Context) error {
	// Start serving in background.
	go s.srv.RunOrDie()

	// Block until the context is canceled or terminated.
	// The following code is used to perform some cleanup tasks when the server shuts down.
	<-ctx.Done()
	slog.Info("Shutting down server...") // Graceful stop server with timeout derived from ctx.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.srv.GracefulStop(ctx)

	slog.Info("Server exited successfully.")

	return nil
}

// NewDB creates and returns a *gorm.DB instance for MySQL.
func (cfg *Config) NewDB() (*gorm.DB, error) {
	slog.Info("Initializing database connection", "type", "mysql")
	db, err := cfg.MySQLOptions.NewDB()
	if err != nil {
		slog.Error("Failed to create database connection", "error", err)
		return nil, err
	}

	// Automatically migrate database schema
	if err := registry.Migrate(db); err != nil {
		slog.Error("Failed to migrate database schema", "error", err)
		return nil, err
	}

	return db, nil
}

// UserRetriever 定义一个用户数据获取器. 用来获取用户信息.
type UserRetriever struct {
	store store.IStore
}

// GetUser 根据用户 ID 获取用户信息.
func (r *UserRetriever) GetUser(ctx context.Context, userID string) (*model.UserM, error) {
	return r.store.User().Get(ctx, where.F("userID", userID))
}

// ProvideDB provides a database instance based on the configuration.
func ProvideDB(cfg *Config) (*gorm.DB, error) {
	return cfg.NewDB()
}

func NewWebServer(serverConfig *ServerConfig, authn authn.Authenticator) (server.Server, error) {
	return serverConfig.NewGinServer(authn)
}

// NewAuthenticator creates a new JWT-based Authenticator using the provided JWT and Redis options.
func NewAuthenticator(jwtOpts *genericoptions.JWTOptions, redisOpts *genericoptions.RedisOptions) (authn.Authenticator, error) {
	// Create a list of options for jwtauthn.
	opts := []jwtauthn.Option{
		// Specify the issuer of the token
		jwtauthn.WithIssuer("art-design-pro-go"),
		// Specify the default expiration time for the token to be issued
		jwtauthn.WithExpired(jwtOpts.Expired),
		// Specify the key to be used when issuing the token
		jwtauthn.WithSigningKey([]byte(jwtOpts.Key)),
		// WithKeyfunc will be used by the Parse methods as a callback function to supply
		// the key for verification.  The function receives the parsed,
		// but unverified Token.  This allows you to use properties in the
		// Header of the token (such as `kid`) to identify which key to use.
		jwtauthn.WithKeyfunc(func(t *jwt.Token) (any, error) {
			// Verify that the signing method is HMAC.
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwtauthn.ErrTokenInvalid
			}
			// Return the signing key.
			return []byte(jwtOpts.Key), nil
		}),
	}

	// Set the signing method based on the provided option.
	var method jwt.SigningMethod
	switch jwtOpts.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}

	opts = append(opts, jwtauthn.WithSigningMethod(method))

	// Create a Redis store for jwtauthn.
	store := jwtredis.NewStore(&jwtredis.Config{
		Addr:      redisOpts.Addr,
		Username:  redisOpts.Username,
		Password:  redisOpts.Password,
		Database:  redisOpts.Database,
		KeyPrefix: "authn_",
	})

	// Create a new jwtauthn instance using the Redis store and options.
	authn := jwtauthn.New(store, opts...)

	return authn, nil
}
