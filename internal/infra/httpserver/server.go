package httpserver

import (
    "context"
    "fragments/internal/core/models"
    "log"
    "log/slog"
    "net/http"
    "os/signal"
    "syscall"
    "time"

    helmet "github.com/danielkov/gin-helmet"
    "github.com/gin-gonic/gin"
)

type Server struct {
    address     string
    Router      *gin.Engine
    requireAuth func(c *gin.Context)
    requireRole func(c *gin.Context, role models.Role)
}

type authRequirer func(c *gin.Context)
type roleRequirer func(c *gin.Context, role models.Role)

func New(address string, requireAuth authRequirer, requireRole roleRequirer) *Server {
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(helmet.Default())
    r.Use(LogError())
    r.Use(LogAccess())
    r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
    r.GET("/version", func(c *gin.Context) { c.JSON(200, gin.H{"message": "123"}) })
    return &Server{address: address, Router: r, requireAuth: requireAuth, requireRole: requireRole}
}

func (s Server) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        s.requireAuth(c)
    }
}

func (s Server) RequireRole(role models.Role) gin.HandlerFunc {
    return func(c *gin.Context) {
        s.requireRole(c, role)
    }
}

func (s Server) Use(middleware gin.HandlerFunc) {
    s.Router.Use(middleware)
}

func (s Server) Run() error {
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    server := &http.Server{Addr: s.address, Handler: s.Router}
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    <-ctx.Done()
    stop()
    slog.Info("shuting down")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        slog.Error("forced shutdown", "err", err)
        return err
    }

    slog.Info("exit")
    return nil
}
