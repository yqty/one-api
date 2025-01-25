package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/yqty/one-api/common"
	"github.com/yqty/one-api/common/client"
	"github.com/yqty/one-api/common/config"
	"github.com/yqty/one-api/common/logger"
	"github.com/yqty/one-api/middleware"
	"github.com/yqty/one-api/model"
	"github.com/yqty/one-api/relay/adaptor/openai"
	"github.com/yqty/one-api/router"
)

func init() {
	common.Init()
	logger.SetupLogger()
	logger.SysLogf("One API %s started", common.Version)

	if os.Getenv("GIN_MODE") != gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	if config.DebugEnabled {
		logger.SysLog("running in debug mode")
	}

	// Initialize SQL Database
	model.InitDB()
	model.InitLogDB()

	var err error
	err = model.CreateRootAccountIfNeed()
	if err != nil {
		logger.FatalLog("database init error: " + err.Error())
	}

	// Initialize Redis
	err = common.InitRedisClient()
	if err != nil {
		logger.FatalLog("failed to initialize Redis: " + err.Error())
	}

	// Initialize options
	model.InitOptionMap()
	logger.SysLog(fmt.Sprintf("using theme %s", config.Theme))
	if common.RedisEnabled {
		config.MemoryCacheEnabled = true
	}
	if config.MemoryCacheEnabled {
		logger.SysLog("memory cache enabled")
		logger.SysLog(fmt.Sprintf("sync frequency: %d seconds", config.SyncFrequency))
		model.InitChannelCache()
	}

	openai.InitTokenEncoders()
	client.Init()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(middleware.RequestId())
	middleware.SetUpLogger(server)

	// Initialize session store
	store := cookie.NewStore([]byte(config.SessionSecret))
	server.Use(sessions.Sessions("session", store))

	router.SetRouter(server, nil)
	server.ServeHTTP(w, r)
}
