package middleware

import "github.com/gin-gonic/gin"

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tiêu đề CORS khớp với Nginx
		//c.Writer.Header().Set("Access-Control-Allow-Origin", "https://isymptom.vercel.app")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization,X-API-Key")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length,Content-Range")

		// Xử lý yêu cầu OPTIONS (preflight)
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Max-Age", "1728000")
			c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
			c.Writer.Header().Set("Content-Length", "0")
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
