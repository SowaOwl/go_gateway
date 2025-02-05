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
	"sync"
)

type LogMiddleware struct {
	db *gorm.DB
}

func NewLogMiddleware(db *gorm.DB) *LogMiddleware {
	return &LogMiddleware{
		db: db,
	}
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

func (l *LogMiddleware) Handle() gin.HandlerFunc {
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

		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()

		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           buf,
		}
		c.Writer = rw

		c.Next()

		go l.writeToDB(c, reqBody, rw)

		bufferPool.Put(buf)
	}
}

func (l *LogMiddleware) writeToDB(c *gin.Context, reqBody []byte, rw *responseWriter) {
	reqHeader, err := json.Marshal(c.Request.Header)
	if err != nil {
		util.SaveErrToDB(err, l.db)
		return
	}

	log := model.ApiLog{
		UserID:        0,
		RequestMethod: c.Request.Method,
		Url:           getScheme(c) + "://" + c.Request.Host + c.Request.URL.String(),
		RequestBody:   string(reqBody),
		RequestHeader: string(reqHeader),
		Ip:            c.ClientIP(),
		ResponseCode:  uint(c.Writer.Status()),
		ResponseBody:  rw.body.String(),
	}

	go func(log model.ApiLog) {
		l.db.Create(&log)
	}(log)
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
