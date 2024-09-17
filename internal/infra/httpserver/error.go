package httpserver

import (
    "fmt"
    "log/slog"

    "github.com/gin-gonic/gin"
)

type ErrorResponse struct {
    Id     string `json:"id" example:"Request Id"`
    Status string `json:"status" example:"Response Status"`
    Error  string `json:"error" example:"Error Message"`
}

func HandleError(c *gin.Context, status int, msg string, err error) {
    id := "-"
    _id, ok := c.Get("requestId")
    if ok && _id != nil {
        id = _id.(string)
    }
    slog.Error(msg, slog.String("id", id), slog.Int("status", status), slog.Any("err", err))
    c.AbortWithStatusJSON(status, ErrorResponse{Id: id, Status: "error", Error: msg})
}

func LogError() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        for i, err := range c.Errors {
            slog.Error(fmt.Sprintf("error %d", i), slog.Any("error", err))
        }
    }
}
