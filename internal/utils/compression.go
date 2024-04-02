package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

var validCompressionContentType = [2]string{"application/json", "text/html"}

type gzipWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (cw gzipWriter) Write(b []byte) (int, error) {
	return cw.Writer.Write(b)
}

func CustomCompression() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isValidContentType := false
		contentType := ctx.Request.Header.Get("Content-Type")
		isNeedCompression := strings.Contains(ctx.Request.Header.Get("Accept-Encoding"), "gzip")

		for _, c := range validCompressionContentType {
			if contentType == c {
				isValidContentType = true
				break
			}
		}

		if isValidContentType && isNeedCompression {
			gz := gzip.NewWriter(ctx.Writer)
			defer func(compressWriter *gzip.Writer) {
				err := compressWriter.Close()
				if err != nil {
					Log.Errorf("Error gz.Close: %s", err)
				}
			}(gz)
			ctx.Header("Content-Encoding", "gzip")
			ctx.Writer = &gzipWriter{ctx.Writer, gz}
		}

		if ctx.Request.Header.Get(`Content-Encoding`) == "gzip" {
			gz, err := gzip.NewReader(ctx.Request.Body)
			if err != nil {
				Log.Errorf("Error NewReader(body): %s", err)
				return
			}
			defer func(gz *gzip.Reader) {
				err := gz.Close()
				if err != nil {
					Log.Errorf("Error gz.Close: %s", err)
				}
			}(gz)

			body, err := io.ReadAll(gz)
			if err != nil {
				log.Fatalf("error: read body: %d", err)
				return
			}

			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
			ctx.Request.ContentLength = int64(len(body))
		}
		ctx.Next()
	}
}
