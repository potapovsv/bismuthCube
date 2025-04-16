🔥 Этап 1: Минимальный работающий XMLA endpoint (2 недели)
Цель: Принимать и отвечать на базовые XMLA Discover-запросы

go
Copy
// main.go - стартовая точка
package main

import (
    "encoding/xml"
    "fmt"
    "log"
    "net/http"
    "strings"
)

func main() {
    http.HandleFunc("/xmla", xmlaHandler)
    log.Println("XMLA server starting on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func xmlaHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var envelope XMLAEnvelope
    if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    response := handleXMLARequest(envelope)
    w.Header().Set("Content-Type", "text/xml")
    xml.NewEncoder(w).Encode(response)
}
📦 Этап 2: Ядро метаданных (3 недели)
Реализуем структуры для хранения метаданных:

go
Copy
// meta/models.go
type Cube struct {
    Name        string
    Dimensions  []Dimension
    Measures    []Measure
    Partitions  []Partition
    StorageMode StorageType // ROLAP, MOLAP, HOLAP
}

type Dimension struct {
    Name       string
    Attributes []Attribute
    Hierarchies []Hierarchy
}

type Measure struct {
    Name      string
    DataType  string
    Aggregator string // SUM, COUNT, etc.
}
Инициализация из TMSL:

go
Copy
func LoadModelFromTMSL(jsonStr string) (*Model, error) {
    var tmslModel TMSLModel
    if err := json.Unmarshal([]byte(jsonStr), &tmslModel); err != nil {
        return nil, err
    }
    // Конвертация TMSL -> внутренняя модель
}
🛠️ Этап 3: MDX парсер и исполнитель (4 недели)
Используем ANTLR для грамматики MDX:

go
Copy
// mdx/parser.go
import (
    "github.com/antlr/antlr4/runtime/Go/antlr"
    "./parser" // сгенерированный парсер
)

type MDXVisitor struct {
    *parser.BaseMDXParserVisitor
}

func ParseMDX(query string) (MDXQuery, error) {
    input := antlr.NewInputStream(query)
    lexer := parser.NewMDXLexer(input)
    stream := antlr.NewCommonTokenStream(lexer, 0)
    p := parser.NewMDXParser(stream)
    
    visitor := &MDXVisitor{}
    return p.Query().Accept(visitor).(MDXQuery), nil
}
🚀 Этап 4: Оптимизация запросов (2 недели)
go
Copy
// engine/optimizer.go
type QueryOptimizer struct {
    Metadata *meta.Model
    Stats    StatisticsCollector
}

func (q *QueryOptimizer) Optimize(query MDXQuery) ExecutionPlan {
    // 1. Применяем правила перезаписи
    rewritten := q.applyRewriteRules(query)
    
    // 2. Генерируем физический план
    plan := q.generatePhysicalPlan(rewritten)
    
    // 3. Выбираем оптимальный порядок выполнения
    return q.optimizeExecutionOrder(plan)
}
📊 Этап 5: Хранилище агрегатов (3 недели)
Реализация агрегатного хранилища:

go
Copy
// storage/aggregate_store.go
type AggregateStore struct {
    db         *sql.DB
    cache      *ristretto.Cache
    aggregator Aggregator
}

func (s *AggregateStore) GetAggregate(
    cube string,
    dimensions []DimensionSlice,
    measures []string,
) (*ResultSet, error) {
    // 1. Проверка кэша
    if agg, found := s.checkCache(cube, dimensions); found {
        return agg, nil
    }
    
    // 2. Проверка предварительно вычисленных агрегатов
    if agg, found := s.checkPrecomputed(cube, dimensions); found {
        return agg, nil
    }
    
    // 3. Вычисление на лету
    return s.computeOnTheFly(cube, dimensions, measures)
}
📅 План разработки (по неделям):
Неделя 1-2: XMLA endpoint + Discover

Неделя 3-5: Ядро метаданных + TMSL загрузка

Неделя 6-9: MDX парсер + базовый исполнитель

Неделя 10-11: Оптимизатор запросов

Неделя 12-14: Хранилище агрегатов

Неделя 15: Интеграционное тестирование