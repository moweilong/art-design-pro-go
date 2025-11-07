# Kratos认证中间件转换为Gin框架使用指南

本文档详细说明如何将Kratos框架的JWT认证中间件转换为Gin框架使用。

## 1. 转换原理

Kratos和Gin的中间件机制有以下主要区别：

| 特性 | Kratos | Gin |
|------|--------|-----|
| 中间件类型 | `middleware.Middleware` | `gin.HandlerFunc` |
| 上下文对象 | `context.Context` | `*gin.Context` |
| 处理流程 | `(ctx context.Context, rq any) (any, error)` | 链式调用，使用`c.Next()`或`c.Abort()` |
| 错误处理 | 返回错误对象 | 使用`c.AbortWithStatusJSON()`等方法直接响应 |

## 2. 实现细节

已创建`authn_jwt.go`文件，实现了Gin版本的JWT认证中间件，主要功能包括：

1. **认证功能**：验证JWT令牌的有效性
2. **上下文注入**：将用户信息注入Gin上下文
3. **错误处理**：提供友好的错误响应
4. **灵活配置**：支持跳过特定路径的认证

### 2.1 核心实现分析

#### Kratos版本（简化）：
```go
func Server(a authn.Authenticator) middleware.Middleware {
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, rq any) (any, error) {
            if tr, ok := transport.FromServerContext(ctx); ok {
                // 从请求头获取令牌
                auths := strings.SplitN(tr.RequestHeader().Get(authorizationKey), " ", 2)
                // 验证令牌
                claims, err := a.ParseClaims(ctx, accessToken)
                // 注入上下文
                ctx = contextx.WithUserID(ctx, claims.Subject)
                // 继续处理请求
                return handler(ctx, rq)
            }
        }
    }
}
```

#### Gin版本实现：
```go
func AuthnJWT(a authn.Authenticator) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从Gin请求头获取令牌
        authHeader := c.GetHeader(authorizationKey)
        // 验证令牌
        claims, err := a.ParseClaims(ctx, accessToken)
        // 注入到Gin上下文和请求上下文
        c.Set("userID", claims.Subject)
        // 继续处理请求链
        c.Next()
    }
}
```

## 3. 主要差异转换

### 3.1 中间件结构转换

- **Kratos**: 使用三层嵌套函数，返回`middleware.Handler`
- **Gin**: 直接返回`gin.HandlerFunc`函数，接收`*gin.Context`参数

### 3.2 请求头获取方式

- **Kratos**: 通过`transport.FromServerContext`获取上下文，再调用`RequestHeader().Get()`
- **Gin**: 直接使用`c.GetHeader()`获取

### 3.3 错误处理方式

- **Kratos**: 返回错误对象，由框架统一处理
- **Gin**: 使用`c.AbortWithStatusJSON()`直接返回HTTP错误响应

### 3.4 上下文数据传递

- **Kratos**: 创建新的上下文并返回
- **Gin**: 
  1. 使用`c.Set()`存储到Gin上下文
  2. 同时更新`c.Request.Context()`以保持与原有逻辑兼容

## 4. 使用方法

### 4.1 基本使用

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/onexstack/onex/pkg/authn"
    "github.com/onexstack/onex/pkg/middleware/gin"
)

func main() {
    // 创建Gin引擎
    r := gin.Default()
    
    // 获取已配置的authn实例
    var authenticator authn.Authenticator = ...
    
    // 使用JWT认证中间件
    r.Use(gin.AuthnJWT(authenticator))
    
    // 路由定义
    r.GET("/api/protected", func(c *gin.Context) {
        // 从上下文获取用户ID
        userID := gin.GetUserID(c)
        c.JSON(200, gin.H{
            "message": "Protected resource",
            "userID":  userID,
        })
    })
    
    r.Run(":8080")
}
```

### 4.2 选择性应用

对于需要跳过认证的路径（如登录接口），可以使用`AuthnJWTSkip`：

```go
func main() {
    r := gin.Default()
    
    // 获取authn实例
    var authenticator authn.Authenticator = ...
    
    // 跳过特定路径的认证
    r.Use(gin.AuthnJWTSkip(authenticator, "/api/login", "/api/register"))
    
    // 公共路由
    r.POST("/api/login", handleLogin)
    r.POST("/api/register", handleRegister)
    
    // 受保护的路由会自动进行认证
    api := r.Group("/api")
    {
        api.GET("/profile", func(c *gin.Context) {
            userID := gin.GetUserID(c)
            // 处理请求...
        })
    }
}
```

### 4.3 与现有authn组件集成

Gin版本的中间件与Kratos版本使用相同的`authn.Authenticator`接口，因此可以直接复用现有的JWT认证器配置：

```go
func createAuthenticator() (authn.Authenticator, error) {
    // 复用现有的JWT认证器配置
    opts := []jwtauthn.Option{
        jwtauthn.WithIssuer("onex-usercenter"),
        jwtauthn.WithExpired(time.Hour * 24),
        jwtauthn.WithSigningKey([]byte("your-secret-key")),
        jwtauthn.WithSigningMethod(jwt.SigningMethodHS256),
    }
    
    // 创建Redis存储（如果需要）
    store := redis.NewStore(&redis.Config{
        Addr: "localhost:6379",
    })
    
    return jwtauthn.New(store, opts...), nil
}
```

## 5. 辅助函数

为了方便使用，提供了以下辅助函数：

- `GetUserID(c *gin.Context) string` - 从上下文获取用户ID
- `GetAccessToken(c *gin.Context) string` - 从上下文获取访问令牌

## 6. 代码优化建议

1. **错误处理增强**：可以添加更细粒度的错误类型处理和日志记录

```go
// 优化版错误处理
if err != nil {
    logger := log.FromContext(ctx)
    logger.Errorw(err, "JWT authentication failed", "path", c.Request.URL.Path)
    
    statusCode := http.StatusUnauthorized
    message := "Unauthorized"
    
    if errors.Is(err, jwtauthn.ErrTokenExpired) {
        statusCode = http.StatusUnauthorized
        message = "Token expired"
    } else if errors.Is(err, jwtauthn.ErrTokenInvalid) {
        statusCode = http.StatusUnauthorized
        message = "Invalid token"
    }
    
    c.AbortWithStatusJSON(statusCode, gin.H{
        "code":    statusCode,
        "message": message,
    })
    return
}
```

2. **性能优化**：对于高频访问的API，可以考虑添加令牌缓存

```go
import "github.com/patrickmn/go-cache"

var tokenCache *cache.Cache

func init() {
    // 5分钟过期，10分钟清理一次
    tokenCache = cache.New(5*time.Minute, 10*time.Minute)
}

func AuthnJWTWithCache(a authn.Authenticator) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader(authorizationKey)
        // ...省略前面的代码...
        
        // 尝试从缓存获取验证结果
        cacheKey := "authn:" + accessToken
        if cachedClaims, found := tokenCache.Get(cacheKey); found {
            claims := cachedClaims.(*jwt.RegisteredClaims)
            // 设置上下文...
            c.Next()
            return
        }
        
        // 缓存未命中，正常验证
        claims, err := a.ParseClaims(ctx, accessToken)
        if err == nil {
            // 缓存验证结果
            tokenCache.Set(cacheKey, claims, cache.DefaultExpiration)
        }
        
        // 设置上下文...
        c.Next()
    }
}
```

## 7. 完整转换示例

下面是将原始Kratos认证中间件完全转换为Gin版本的完整流程：

1. 创建`authn_jwt.go`文件，实现Gin版中间件
2. 创建辅助函数简化上下文访问
3. 在Gin路由中应用中间件
4. 在处理器中获取认证信息

这样就可以在不改变底层认证逻辑的情况下，将Kratos的认证中间件无缝转换为Gin框架使用。