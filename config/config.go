package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

//type LoggingConfig struct {
//	Enabled bool   `yaml:"enabled"`
//	File    string `yaml:"file"`
//	Level   string `yaml:"level"`
//}

type XmlaConfig struct {
	DataSource struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		URL         string `yaml:"url"`
	} `yaml:"datasource"`

	Catalog struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	} `yaml:"catalog"`

	Server struct {
		Port        int    `yaml:"port"`
		Version     string `yaml:"version"`
		ProductName string `yaml:"product name"`
		//isShowConfig bool   `yaml:"isShowConfig"`
	} `yaml:"server"`
	Logging struct {
		Enabled bool   `yaml:"enabled"`
		File    string `yaml:"file"`
		Level   string `yaml:"level"`
	} `yaml:"logging"`
}

var (
	instance *XmlaConfig
	once     sync.Once
)

// GetConfig возвращает синглтон конфигурации
func GetConfig() *XmlaConfig {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}
func findConfigFile() string {
	// Возможные места расположения конфига (относительно корня проекта)
	tryPaths := []string{
		"config.yml",                                // 1. Рядом с исполняемым файлом
		filepath.Join("config", "config.yml"),       // 2. В папке config/
		filepath.Join("..", "config.yml"),           // 3. На уровень выше (для bin/)
		filepath.Join("..", "config", "config.yml"), // 4. На уровень выше в config/
	}

	// Добавляем абсолютные пути на основе расположения этого файла
	if _, filename, _, ok := runtime.Caller(0); ok {
		baseDir := filepath.Dir(filename)
		tryPaths = append(tryPaths,
			filepath.Join(baseDir, "config.yml"),
			filepath.Join(baseDir, "..", "config.yml"),
			filepath.Join(baseDir, "..", "config", "config.yml"),
		)
	}

	// Проверяем все возможные пути
	for _, path := range tryPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	log.Fatal("Config file not found in:\n" + strings.Join(tryPaths, "\n"))
	return ""
}
func loadConfig() *XmlaConfig {
	cfg := XmlaConfig{}
	configPath := findConfigFile()
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	// Устанавливаем значения по умолчанию
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}

	return &cfg
}
func (c *XmlaConfig) JSON() string {
	jsonData, _ := json.MarshalIndent(c, "", "  ")
	return string(jsonData)
}
func (c *XmlaConfig) String() string {
	var sb strings.Builder
	//if GetConfig().Server.isShowConfig == false {
	//	sb.WriteString("\n⚙️ Configuration:isShowConfig=false \n")
	//	return sb.String()
	//}

	val := reflect.ValueOf(c).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Получаем описание из тега
		desc := fieldType.Tag.Get("desc")
		if desc == "" {
			desc = fieldType.Name
		}

		sb.WriteString(fmt.Sprintf("\n%s:\n", desc))

		if field.Kind() == reflect.Struct {
			for j := 0; j < field.NumField(); j++ {
				nestedField := field.Field(j)
				nestedType := field.Type().Field(j)

				nestedDesc := nestedType.Tag.Get("desc")
				if nestedDesc == "" {
					nestedDesc = nestedType.Name
				}

				sb.WriteString(fmt.Sprintf("  - %s: %v\n", nestedDesc, nestedField.Interface()))
			}
		}
	}

	return sb.String()
}
