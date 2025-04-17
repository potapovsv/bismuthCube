package main

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/potapovsv/bismuthCube/config"
	"github.com/potapovsv/bismuthCube/protocols/xmla"
	"net/http"
	"strconv"
)

func main() {
	_ = config.GetConfig()
	r := gin.Default()

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
