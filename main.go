package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Embed the HTML template into the binary at compile time
//go:embed index.html
var htmlTemplate string

type Settings struct {
	CompanyName string `yaml:"company_name"`
	LogoURL     string `yaml:"logo_url"`
	FooterText  string `yaml:"footer_text"`
}

type App struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
	Icon        string `yaml:"icon"`
}

type ConfigFile struct {
	Environment string `yaml:"environment"`
	Category    string `yaml:"category"`
	Apps        []App  `yaml:"apps"`
}

type PageData struct {
	Settings Settings
	Data     map[string]map[string][]App
}

func main() {
	http.HandleFunc("/", handleDashboard)
	log.Println("Dashboard server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	// 1. Set default values in case settings.yaml is missing
	settings := Settings{
		CompanyName: "My Company",
		FooterText:  "© 2026 My Company. All rights reserved.",
	}

	// Try to load settings.yaml
	settingsPath := "/config/settings.yaml"
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		settingsPath = "./config/settings.yaml" // Local fallback for development
	}

	if data, err := os.ReadFile(settingsPath); err == nil {
		if err := yaml.Unmarshal(data, &settings); err != nil {
			log.Printf("Error parsing settings.yaml: %v", err)
		}
	}

	// 2. Load application configuration files
	files, err := filepath.Glob("/config/*.yaml")
	if err != nil || len(files) == 0 {
		files, _ = filepath.Glob("./config/*.yaml")
	}

	// Structure: [Environment][Category] -> List of Apps
	dashboardData := make(map[string]map[string][]App)

	for _, file := range files {
		// Skip the global settings file
		if filepath.Base(file) == "settings.yaml" {
			continue
		}

		data, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error reading file %s: %v", file, err)
			continue
		}

		var config ConfigFile
		if err := yaml.Unmarshal(data, &config); err != nil {
			log.Printf("Error parsing file %s: %v", file, err)
			continue
		}

		if dashboardData[config.Environment] == nil {
			dashboardData[config.Environment] = make(map[string][]App)
		}
		dashboardData[config.Environment][config.Category] = append(dashboardData[config.Environment][config.Category], config.Apps...)
	}

	tmpl := template.Must(template.New("index").Parse(htmlTemplate))
	tmpl.Execute(w, PageData{Settings: settings, Data: dashboardData})
}
