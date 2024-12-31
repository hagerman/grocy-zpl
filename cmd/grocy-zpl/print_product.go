package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/OpenPrinting/goipp"
	"github.com/hagerman/grocy-zpl/internal/funcs"
	"github.com/hagerman/grocy-zpl/internal/services"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

// ProductWebhook represents the expected JSON payload structure
type ProductWebhook struct {
	Name       string `json:"product"`
	Barcode    string `json:"grocycode"`
	DueDateRaw string `json:"due_date"`
	DueDate    time.Time
}

type Product struct {
	Webhook     ProductWebhook
	ServiceCall services.ProductResponse
	MediaReady  string
}

// handler for POST requests to /print/product
func (app *application) printProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var webhook ProductWebhook
	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if webhook.DueDateRaw != "" {
		d, err := time.Parse("DD: 2006-01-02", webhook.DueDateRaw)
		if err == nil { // Note this is backwards, if no error, we want to persist the value
			webhook.DueDate = d
		}
	}

	var serviceCall services.ProductResponse
	var mediaReady string
	if app.grocyAPIURL != "" {
		serviceCall, err = services.GetProductByBarcode(app.grocyAPIURL, app.grocyAPIKey, webhook.Barcode)
		if err != nil {
			http.Error(w, "Unable to get product by barcode", http.StatusInternalServerError)
			return
		}
	}

	mediaReady, err = getMediaReadyAttr(app.printerURL)
	if err != nil {
		http.Error(w, "Unable to get media-ready attribute from printer", http.StatusInternalServerError)
	}

	product := Product{
		Webhook:     webhook,
		ServiceCall: serviceCall,
		MediaReady:  mediaReady,
	}

	log.Printf("Printing product: %+v\n", webhook)
	err = app.printProduct(product)
	if err != nil {
		log.Printf("Error printing product: %+v\n", err)
		http.Error(w, fmt.Sprintf("Error printing product: %+v\n", err), http.StatusBadRequest)
		return
	}

	log.Printf("Successfully printed: %+v\n", webhook)

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product printed successfully"))
}

func (app *application) printProduct(product Product) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Construct the path to the file in the root directory
	templatePath := filepath.Join(cwd, app.templatePath)

	// Parse the template
	tmpl := template.Must(template.New("").Funcs(funcs.TemplateFuncs).ParseFiles(templatePath))

	// Replace contents using the template
	var fileData bytes.Buffer
	err = tmpl.ExecuteTemplate(&fileData, filepath.Base(templatePath), product)
	if err != nil {
		return err
	}

	// Create an IPP Print-Job request
	req := goipp.NewRequest(goipp.DefaultVersion, goipp.OpPrintJob, 1)
	req.Operation.Add(goipp.MakeAttribute("attributes-charset", goipp.TagCharset, goipp.String("utf-8")))
	req.Operation.Add(goipp.MakeAttribute("attributes-natural-language", goipp.TagLanguage, goipp.String("en-US")))
	req.Operation.Add(goipp.MakeAttribute("printer-uri", goipp.TagURI, goipp.String(app.printerURL)))
	req.Operation.Add(goipp.MakeAttribute("document-format", goipp.TagMimeType, goipp.String("application/vnd.zebra-zpl")))

	req.Job.Add(makeAttrCollection("media-col",
		goipp.MakeAttribute("media-left-margin",
			goipp.TagInteger, goipp.Integer(0)),
		goipp.MakeAttribute("media-right-margin",
			goipp.TagInteger, goipp.Integer(0)),
		goipp.MakeAttribute("media-top-margin",
			goipp.TagInteger, goipp.Integer(0)),
		goipp.MakeAttribute("media-bottom-margin",
			goipp.TagInteger, goipp.Integer(0)),
	))

	// Encode the IPP request into a buffer
	var buffer bytes.Buffer
	err = req.Encode(&buffer)
	if err != nil {
		return fmt.Errorf("Failed to encode IPP request: %v", err)
	}

	// Append the ZPL data to the buffer
	buffer.Write(fileData.Bytes())

	// Send the IPP request to the printer using HTTP
	httpReq, err := http.NewRequest("POST", app.printerURL, &buffer)
	if err != nil {
		return fmt.Errorf("Failed to create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/ipp")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("Failed to send request: %v", err)
	}
	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Failed to close response body: %v", err)
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send request, received %s", resp.Status)
	}

	log.Printf("ZPL job submitted successfully!")
	return nil
}
