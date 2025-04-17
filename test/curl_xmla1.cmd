curl -X POST http://localhost:8080/xmla -H "Content-Type: text/xml" -d '
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <Discover xmlns="urn:schemas-microsoft-com:xml-analysis">
      <RequestType>DISCOVER_DATASOURCES</RequestType>
    </Discover>
  </soap:Body>
</soap:Envelope>'