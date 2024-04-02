package utils

import (
	"compress/gzip"
	"io"
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

func CustomCompression(ctx *gin.Context) {

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

	isNeedCompression := strings.Contains(ctx.Request.Header.Get("Accept-Encoding"), "gzip")
	if !isNeedCompression {
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

	cw := &gzipWriter{Writer: gz, ResponseWriter: ctx.Writer}
	ctx.Writer = cw
	ctx.Next()
}
