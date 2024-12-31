package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type application struct {
	printerURL   string
	templatePath string
	grocyAPIURL  string
	grocyAPIKey  string
}

func main() {
	// Get the port from the environment variable or use the default value
	port := getEnv("PORT", "8000")
	printerURL := getEnv("PRINTER_URL", "http://1.1.1.1:631/ipp/print")
	templatePath := getEnv("TEMPLATE_PATH", "assets/templates/product.zpl")
	grocyAPIURL := getEnv("GROCY_API_URL", "")
	grocyAPIKey := getEnv("GROCY_API_KEY", "")

	app := &application{
		printerURL:   printerURL,
		templatePath: templatePath,
		grocyAPIURL:  grocyAPIURL,
		grocyAPIKey:  grocyAPIKey,
	}

	http.HandleFunc("/print/product", app.printProductHandler)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on port %s. Printing to URL %s...", port, printerURL)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
