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
}

func main() {
	// Get the port from the environment variable or use the default value
	port := getEnv("PORT", "8000")
	printerURL := getEnv("PRINTER_URL", "http://1.1.1.1:631/ipp/print")
	templatePath := getEnv("TEMPLATE_PATH", "template.zpl")

	app := &application{
		printerURL:   printerURL,
		templatePath: templatePath,
	}

	http.HandleFunc("/print/product", app.printProductHandler)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on port %s...", port)
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
