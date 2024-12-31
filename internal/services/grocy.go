package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type UserFields map[string]interface{}

type ProductResponse struct {
	Product struct {
		ID                  int     `json:"id"`
		Name                string  `json:"name"`
		Description         *string `json:"description"`
		ProductGroupID      int     `json:"product_group_id"`
		Active              int     `json:"active"`
		LocationID          int     `json:"location_id"`
		ShoppingLocationID  int     `json:"shopping_location_id"`
		QuIDPurchase        int     `json:"qu_id_purchase"`
		QuIDStock           int     `json:"qu_id_stock"`
		MinStockAmount      int     `json:"min_stock_amount"`
		PictureFileName     string  `json:"picture_file_name"`
		RowCreatedTimestamp string  `json:"row_created_timestamp"`
	} `json:"product"`
	ProductBarcodes []struct {
		ID        int      `json:"id"`
		ProductID int      `json:"product_id"`
		Barcode   string   `json:"barcode"`
		QuID      int      `json:"qu_id"`
		Amount    int      `json:"amount"`
		LastPrice *float64 `json:"last_price"`
	} `json:"product_barcodes"`
	StockAmount int     `json:"stock_amount"`
	StockValue  float64 `json:"stock_value"`
	Location    struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		IsFreezer int    `json:"is_freezer"`
		Active    int    `json:"active"`
	} `json:"location"`
	LastPrice              float64 `json:"last_price"`
	CurrentPrice           float64 `json:"current_price"`
	ProductUserFields      UserFields
	ProductGroupUserFields UserFields
}

func GetProductByBarcode(apiURL, apiKey, barcode string) (ProductResponse, error) {
	client := &http.Client{}

	u := fmt.Sprintf("%s/stock/products/by-barcode/%s", apiURL, url.QueryEscape(barcode))
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return ProductResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("GROCY-API-KEY", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return ProductResponse{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return ProductResponse{}, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var product ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return ProductResponse{}, fmt.Errorf("error decoding response: %v", err)
	}

	product.ProductUserFields, _ = getUserFields(apiKey, apiURL, "products", product.Product.ID)
	product.ProductGroupUserFields, _ = getUserFields(apiKey, apiURL, "product_groups", product.Product.ProductGroupID)

	return product, nil
}

func getUserFields(apiKey, apiURL, entity string, productID int) (UserFields, error) {
	client := &http.Client{}

	u := fmt.Sprintf("%s/userfields/%s/%d", apiURL, entity, productID)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("GROCY-API-KEY", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userFields UserFields
	if err := json.NewDecoder(resp.Body).Decode(&userFields); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return userFields, nil
}
