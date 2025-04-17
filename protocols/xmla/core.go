package xmla

import (
	"encoding/xml"
	"fmt"
	"github.com/potapovsv/bismuthCube/config"
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
	cfg := config.GetConfig()

	switch requestType {
	case "DISCOVER_DATASOURCES":
		return []byte(fmt.Sprintf(`
			<return xmlns="urn:schemas-microsoft-com:xml-analysis">
				<row>
					<DataSourceName>%s</DataSourceName>
					<DataSourceDescription>%s</DataSourceDescription>
					<URL>%s</URL>
				</row>
			</return>`,
			cfg.DataSource.Name,
			cfg.DataSource.Description,
			cfg.DataSource.URL)), nil
	case "DBSCHEMA_CATALOGS":
		return []byte(fmt.Sprintf(`
			<return xmlns="urn:schemas-microsoft-com:xml-analysis">
				<row>
					<CATALOG_NAME>%s</CATALOG_NAME>
					<DESCRIPTION>%s</DESCRIPTION>
				</row>
			</return>`,
			cfg.Catalog.Name,
			cfg.Catalog.Description)), nil
	default:
		return nil, fmt.Errorf("unsupported request type: %s", requestType)
	}
}
