package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"time"

	"fruit_shop/handler"
	"fruit_shop/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/willf/pad"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logging is a middleware function that logs each request
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path

		reg := regexp.MustCompile("/v1/user|/login")
		if !reg.MatchString(path) {
			return
		}

		// Skip the health check requests
		if path == "/sd/health" || path == "/sd/ram" || path == "/sd/cpu" || path == "/sd/disk" {
			return
		}

		// Read the body content
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// The basic information
		method := c.Request.Method
		ip := c.ClientIP()

		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = blw
		c.Next()

		// Calculates the latency
		end := time.Now().UTC()
		latency := end.Sub(start)

		code, message := -1, ""

		// Get code and message
		var response handler.Response
		if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
			log.Errorf(err, "response body can not unmarshal to model.Response struct, body: `%s`", blw.body.Bytes())
			code = errno.InternalServerError.Code
			message = err.Error()
		} else {
			code = response.Code
			message = response.Message
		}

		log.Infof("%-13s | %-12s | %s %s | {code: %d, message: %s}", latency, ip, pad.Right(method, 5, ""), path, code, message)
	}
}
