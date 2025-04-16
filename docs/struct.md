bismuthCube/
├── core/
│   ├── olap/               # Ядро OLAP-функциональности
│   │   ├── cube.go         # Логика кубов
│   │   ├── dimension.go    # Работа с измерениями
│   │   └── measure.go      # Агрегатные функции
│   └── query/              # Обработка запросов
│       ├── mdx/            # MDX исполнитель
│       └── sql/            # SQL транслятор
├── drivers/
│   ├── clickhouse/         # ClickHouse специфика
│   │   ├── reader.go       # Чтение данных
│   │   ├── writer.go       # Запись агрегатов
│   │   └── optimize.go     # ClickHouse-специфичные оптимизации
│   └── metastore/          # Хранение метаданных
│       ├── postgres/       # (опционально)
│       └── file/           # Для разработки
├── protocols/
│   ├── xmla/               # XML for Analysis
│   ├── tmsl/               # Tabular Model Scripting
│   └── rest/               # REST API
├── services/
│   ├── query/              # Сервис запросов
│   ├── cache/             # Многоуровневый кэш
│   └── monitoring/        # Метрики и логи
├── bin/
│   ├── bismuthd            # Сервер
│   └── bcli               # CLI утилита (bismuthCLI)
└── config/
    ├── clickhouse.yaml    # CH connection config
    └── server.yaml        # Настройки сервера