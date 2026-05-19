package api

import "github.com/cloudwego/hertz/pkg/app"

// writeJSON 写入 JSON 响应。
func writeJSON(c *app.RequestContext, code int, obj interface{}) {
	c.JSON(code, obj)
}

// writeError 写入错误响应。
func writeError(c *app.RequestContext, code int, msg string) {
	c.JSON(code, map[string]string{"error": msg})
}
