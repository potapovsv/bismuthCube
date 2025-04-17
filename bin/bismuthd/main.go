package main

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/potapovsv/bismuthCube/config"
	"github.com/potapovsv/bismuthCube/core/logger"
	"github.com/potapovsv/bismuthCube/protocols/xmla"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	_ = config.GetConfig()
	log := logger.Init(config.GetConfig().Logging.File != "")
	r := gin.Default()
	if config.GetConfig().Logging.File != "" {
		log.Printf("📁 Logging to file: %s", config.GetConfig().Logging.File)
	}

	log.Printf("🚀 Starting BismuthCube server on port %d", config.GetConfig().Server.Port)
	log.Printf("🚀 Version BismuthCube:%s", config.GetConfig().Server.Version)
	log.Printf("Config:\n%s", config.GetConfig().JSON())
	log.Printf("🔗 XMLA endpoint: %s", config.GetConfig().DataSource.URL)

	// Обработка graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		log.Printf("🛑 Server shutting down...")
		os.Exit(0)
	}()
	r.POST("/xmla", func(c *gin.Context) {
		// Парсинг SOAP
		var env xmla.Envelope
		if err := xml.NewDecoder(c.Request.Body).Decode(&env); err != nil {
			c.XML(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Обработка запроса
		response, err := xmla.HandleDiscover(env.Body.Discover.RequestType)
		if err != nil {
			c.XML(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Формируем SOAP-ответ
		c.Header("Content-Type", "text/xml")
		c.String(http.StatusOK, `
            <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
                <soap:Body>
                    %s
                </soap:Body>
            </soap:Envelope>
        `, string(response))
	})

	err := r.Run(":" + strconv.Itoa(config.GetConfig().Server.Port))
	if err != nil {
		return
	}

}
