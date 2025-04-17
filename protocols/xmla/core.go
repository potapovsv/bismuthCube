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
