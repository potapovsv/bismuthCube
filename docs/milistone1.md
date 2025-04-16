
1. Создаем базовую структуру проекта
bash
Copy
# Инициализация модуля
mkdir -p bismuthCube/{protocols/xmla,cmd/bismuthd}
cd bismuthCube
go mod init github.com/yourname/bismuthCube
go get github.com/gin-gonic/gin
2. Реализуем XMLA обработчик (protocols/xmla/core.go)
go
Copy
package xmla

import (
    "encoding/xml"
    "fmt"
)

// SOAP структуры
type Envelope struct {
    XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
    Body    Body     `xml:"Body"`
}

type Body struct {
    Discover *DiscoverRequest `xml:"Discover"`
}

type DiscoverRequest struct {
    XMLName     xml.Name `xml:"urn:schemas-microsoft-com:xml-analysis Discover"`
    RequestType string   `xml:"RequestType"`
}

// Обработчик Discover
func HandleDiscover(requestType string) ([]byte, error) {
    switch requestType {
    case "DISCOVER_DATASOURCES":
        return []byte(`
            <return xmlns="urn:schemas-microsoft-com:xml-analysis">
                <row>
                    <DataSourceName>BismuthCube</DataSourceName>
                    <DataSourceDescription>OLAP server</DataSourceDescription>
                    <URL>http://localhost:8080/xmla</URL>
                </row>
            </return>`), nil
    default:
        return nil, fmt.Errorf("unsupported request type: %s", requestType)
    }
}
3. HTTP сервер (cmd/bismuthd/main.go)
go
Copy
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/yourname/bismuthCube/protocols/xmla"
    "encoding/xml"
    "net/http"
)

func main() {
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

    r.Run(":8080")
}
4. Тестируем сервер
Запускаем сервер:

bash
Copy
go run cmd/bismuthd/main.go
Отправляем тестовый запрос:

bash
Copy
curl -X POST http://localhost:8080/xmla -H "Content-Type: text/xml" -d '
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <Discover xmlns="urn:schemas-microsoft-com:xml-analysis">
      <RequestType>DISCOVER_DATASOURCES</RequestType>
    </Discover>
  </soap:Body>
</soap:Envelope>'
Ожидаемый ответ:

xml
Copy
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <return xmlns="urn:schemas-microsoft-com:xml-analysis">
      <row>
        <DataSourceName>BismuthCube</DataSourceName>
        <DataSourceDescription>OLAP server</DataSourceDescription>
        <URL>http://localhost:8080/xmla</URL>
      </row>
    </return>
  </soap:Body>
</soap:Envelope>
Run HTML
5. Дорабатываем функционал (день 3-4)
Добавим поддержку DBSCHEMA_CATALOGS:

go
Copy
// protocols/xmla/core.go
func HandleDiscover(requestType string) ([]byte, error) {
    switch requestType {
    case "DISCOVER_DATASOURCES":
        return []byte(`...`), nil // предыдущий код
        
    case "DBSCHEMA_CATALOGS":
        return []byte(`
            <return xmlns="urn:schemas-microsoft-com:xml-analysis">
                <row>
                    <CATALOG_NAME>BismuthCube</CATALOG_NAME>
                    <DESCRIPTION>Main catalog</DESCRIPTION>
                </row>
            </return>`), nil
    default:
        return nil, fmt.Errorf("unsupported request type: %s", requestType)
    }
}
6. Добавляем валидацию (день 5-6)
go
Copy
// protocols/xmla/validator.go
package xmla

import "errors"

func (env *Envelope) Validate() error {
    if env.Body.Discover == nil {
        return errors.New("discover request is required")
    }
    if env.Body.Discover.RequestType == "" {
        return errors.New("request type is empty")
    }
    return nil
}

// Обновляем обработчик в main.go:
if err := env.Validate(); err != nil {
    c.XML(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
7. Запускаем финальное тестирование (день 7)
Проверяем все сценарии:

bash
Copy
# 1. Discover DataSources
curl -X POST http://localhost:8080/xmla -d '...'

# 2. Get Catalogs
curl -X POST http://localhost:8080/xmla -d '
<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
  <Body>
    <Discover xmlns="urn:schemas-microsoft-com:xml-analysis">
      <RequestType>DBSCHEMA_CATALOGS</RequestType>
    </Discover>
  </Body>
</Envelope>'

# 3. Невалидный запрос
curl -X POST http://localhost:8080/xmla -d '
<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
  <Body>
    <Discover xmlns="urn:schemas-microsoft-com:xml-analysis">
      <!-- Пропущен RequestType -->
    </Discover>
  </Body>
</Envelope>'
Что дальше?
День 8-9: Добавляем поддержку MDSCHEMA_CUBES

День 10-11: Реализуем кэширование метаданных

День 12-14: Настраиваем логирование и мониторинг
