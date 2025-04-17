package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

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
	} `yaml:"server"`
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
