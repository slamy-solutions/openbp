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

func MetadataFromRequestContext(ctx *gin.Context) *Metadata {
	return &Metadata{
		RealAddress:   ctx.ClientIP(),
		Host:          ctx.GetHeader("Host"),
		XForwardedFor: ctx.GetHeader("X-Forwarded-For"),
		UserAgent:     ctx.GetHeader("User-Agent"),
	}
}

func (m *Metadata) ToJSONString() string {
	data, _ := json.Marshal(m)
	return string(data)
}
