package middlewares

import (
	"bytes"
	"encoding/json"
	"gateway/database/model"
	"gateway/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type LogMiddleware struct {
	db *gorm.DB
}

func NewLogMiddleware(db *gorm.DB) *LogMiddleware {
	return &LogMiddleware{
		db: db,
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

func (l *LogMiddleware) LogRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody []byte

		if c.Request.Header.Get("Content-Type") == "application/json" {
			var err error
			reqBody, err = c.GetRawData()
			if err != nil {
				util.SendError(c, http.StatusInternalServerError, err.Error(), "")
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer([]byte{}),
		}
		c.Writer = rw
		c.Next()

		go l.writeToDB(c, reqBody, rw)
	}
}

func (l *LogMiddleware) writeToDB(c *gin.Context, reqBody []byte, rw *responseWriter) {
	reqHeader, err := json.Marshal(c.Request.Header)
	if err != nil {
		util.SendError(c, http.StatusInternalServerError, err.Error(), "")
	}

	log := model.ApiLog{
		UserID:        0,
		RequestMethod: c.Request.Method,
		Url:           getScheme(c) + "://" + c.Request.Host + c.Request.URL.String(),
		RequestBody:   string(reqBody),
		RequestHeader: string(reqHeader),
		Ip:            c.Request.RemoteAddr,
		ResponseCode:  uint(c.Writer.Status()),
		ResponseBody:  rw.body.String(),
	}

	l.db.Create(&log)
}

func getScheme(c *gin.Context) string {
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		return proto
	}
	if c.Request.TLS != nil {
		return "https"
	}
	return "http"
}
