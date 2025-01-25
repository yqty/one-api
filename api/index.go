package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/yqty/one-api/router"
)

var Handler http.Handler

func init() {
    gin.SetMode(gin.ReleaseMode)
    r := router.SetupRouter()
    Handler = r
}
