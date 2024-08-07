// Package middleware consist methods for parse http request
package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

// Write compression writer
func (cw gzipWriter) Write(b []byte) (int, error) {
	return cw.Writer.Write(b)
}

// Compression - middleware for compress and decompress zip data
func Compression() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Header.Get("Content-Type") == "application/json" ||
			ctx.Request.Header.Get("Content-Type") == "text/html" {
			if strings.Contains(ctx.Request.Header.Get("Accept-Encoding"), "gzip") {
				compressWriter := gzip.NewWriter(ctx.Writer)
				defer func(compressWriter *gzip.Writer) {
					err := compressWriter.Close()
					if err != nil {
						Log.Error("error: compressWriter close: %d", err)
					}
				}(compressWriter)
				ctx.Header("Content-Encoding", "gzip")
				ctx.Writer = &gzipWriter{ctx.Writer, compressWriter}
			}
		}

		if strings.Contains(ctx.Request.Header.Get("Content-Encoding"), "gzip") {
			compressReader, err := gzip.NewReader(ctx.Request.Body)
			if err != nil {
				Log.Error("error: new reader: %d", err)
				return
			}
			defer func(compressReader *gzip.Reader) {
				errClose := compressReader.Close()
				if errClose != nil {
					Log.Error("error: close compress: %d", errClose)
				}
			}(compressReader)

			body, err := io.ReadAll(compressReader)
			if err != nil {
				Log.Error("error: read body: %d", err)
				return
			}

			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
			ctx.Request.ContentLength = int64(len(body))
		}
		ctx.Next()
	}
}
