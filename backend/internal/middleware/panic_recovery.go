package middleware

import (
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "runtime/debug"
)

func PanicRecovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic recovered: %v\n%s", err, debug.Stack())
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal Server Error",
                })
            }
        }()
        c.Next()
    }
}