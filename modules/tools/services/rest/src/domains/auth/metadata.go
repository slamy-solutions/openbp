package auth

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Metadata struct {
	RealAddress string `json:"RealAddress"`

	Host          string `json:"Host"`
	XForwardedFor string `json:"XForwardedFor"`
	UserAgent     string `json:"UserAgent"`
}

func maxHeaderSizeString(s string, max int) string {
	if len(s) > max {
		r := 0
		for i := range s {
			r++
			if r > max {
				return s[:i]
			}
		}
	}
	return s
}

func MetadataFromRequestContext(ctx *gin.Context) *Metadata {
	return &Metadata{
		RealAddress:   ctx.ClientIP(),
		Host:          maxHeaderSizeString(ctx.GetHeader("Host"), 128),
		XForwardedFor: maxHeaderSizeString(ctx.GetHeader("X-Forwarded-For"), 128),
		UserAgent:     maxHeaderSizeString(ctx.GetHeader("User-Agent"), 128),
	}
}

func (m *Metadata) ToJSONString() string {
	data, _ := json.Marshal(m)
	return string(data)
}
