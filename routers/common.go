package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/middlepartedhairstyle/HiWe/utils"
	"io"
	"os"
)

type Server struct {
	Host string
	Port string
}

func NewServer(host string, port string) *Server {
	return &Server{
		Host: host,
		Port: port,
	}
}

func init() {
	_, err := utils.CreateFilePath("images/profile_photo")
	if err != nil {
		panic(err)
	}
}

func (s *Server) Run() {
	gin.SetMode(gin.DebugMode)
	logFile, _ := utils.CreateFile("logs", "gin.log")
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	router := gin.Default()
	Routers(router)

	err := router.Run(s.Host + ":" + s.Port)
	if err != nil {
		return
	}
}
