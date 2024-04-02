package utils

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"strings"
)

var validCompressionContentType = [2]string{"application/json", "text/html"}

type gzipWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (cw gzipWriter) Write(b []byte) (int, error) {
	return cw.buf.Write(b)
}

func CustomCompression(ctx *gin.Context) {
	isValidContentType := false
	contentType := ctx.Request.Header.Get("Content-Type")
	isNeedCompression := strings.Contains(ctx.Request.Header.Get("Accept-Encoding"), "gzip")

	for _, c := range validCompressionContentType {
		if contentType == c {
			isValidContentType = true
			break
		}
	}

	if ctx.Request.Header.Get(`Content-Encoding`) == "gzip" {
		gz, err := gzip.NewReader(ctx.Request.Body)
		if err != nil {
			Sugar.Errorf("Error NewReader(body): %s", err)
			return
		}
		ctx.Request.Body = gz
		defer func(gz *gzip.Reader) {
			err := gz.Close()
			if err != nil {
				Sugar.Errorf("Error gz.Close: %s", err)
			}
		}(gz)
	}

	if !isNeedCompression || !isValidContentType {
		ctx.Next()
		return
	}

	ctx.Writer.Header().Set("Content-Encoding", "gzip")
	gz, err := gzip.NewWriterLevel(ctx.Writer, gzip.BestSpeed)
	if err != nil {
		Sugar.Errorf("Error gzip compression: %s", err)
		return
	}
	defer func(gz *gzip.Writer) {
		err := gz.Close()
		if err != nil {
			Sugar.Errorf("Error gz close: %s", err)
		}
	}(gz)

	cw := &gzipWriter{buf: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
	ctx.Writer = cw
	ctx.Next()
}