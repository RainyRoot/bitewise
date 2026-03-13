package nutrition

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// FoodItem represents a food product with nutritional information.
type FoodItem struct {
	Barcode     string  `json:"barcode"`
	Name        string  `json:"name"`
	Brand       string  `json:"brand,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
	Calories    float64 `json:"calories_per_100g"`
	ProteinG    float64 `json:"protein_g_per_100g"`
	CarbsG      float64 `json:"carbs_g_per_100g"`
	FatG        float64 `json:"fat_g_per_100g"`
	FiberG      float64 `json:"fiber_g_per_100g"`
	ServingSize string  `json:"serving_size,omitempty"`
}

// NutritionProvider looks up food items by barcode or search query.
type NutritionProvider interface {
	LookupBarcode(ctx context.Context, code string) (*FoodItem, error)
	SearchFood(ctx context.Context, query string) ([]FoodItem, error)
}

// OpenFoodFactsProvider implements NutritionProvider using the OpenFoodFacts API.
type OpenFoodFactsProvider struct {
	client  *http.Client
	baseURL string
}

// NewOpenFoodFactsProvider creates a new OpenFoodFacts provider.
func NewOpenFoodFactsProvider() *OpenFoodFactsProvider {
	return &OpenFoodFactsProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://world.openfoodfacts.org",
	}
}

type offProductResponse struct {
	Status  int        `json:"status"`
	Product offProduct `json:"product"`
}

type offProduct struct {
	ProductName string        `json:"product_name"`
	Brands      string        `json:"brands"`
	Code        string        `json:"code"`
	ImageURL    string        `json:"image_url"`
	Nutriments  offNutriments `json:"nutriments"`
	ServingSize string        `json:"serving_size"`
}

type offNutriments struct {
	EnergyKcal100g float64 `json:"energy-kcal_100g"`
	Proteins100g   float64 `json:"proteins_100g"`
	Carbs100g      float64 `json:"carbohydrates_100g"`
	Fat100g        float64 `json:"fat_100g"`
	Fiber100g      float64 `json:"fiber_100g"`
}

type offSearchResponse struct {
	Products []offProduct `json:"products"`
	Count    int          `json:"count"`
}

func (p *OpenFoodFactsProvider) LookupBarcode(ctx context.Context, code string) (*FoodItem, error) {
	reqURL := fmt.Sprintf("%s/api/v2/product/%s.json", p.baseURL, url.PathEscape(code))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "BiteWise/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching product: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openfoodfacts returned status %d", resp.StatusCode)
	}

	var result offProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if result.Status == 0 || result.Product.ProductName == "" {
		return nil, fmt.Errorf("product not found")
	}

	return productToFoodItem(result.Product), nil
}

func (p *OpenFoodFactsProvider) SearchFood(ctx context.Context, query string) ([]FoodItem, error) {
	reqURL := fmt.Sprintf("%s/cgi/search.pl?search_terms=%s&search_simple=1&action=process&json=1&page_size=20",
		p.baseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "BiteWise/1.0")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("searching food: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openfoodfacts returned status %d", resp.StatusCode)
	}

	var result offSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	items := make([]FoodItem, 0, len(result.Products))
	for _, prod := range result.Products {
		if prod.ProductName == "" {
			continue
		}
		items = append(items, *productToFoodItem(prod))
	}

	return items, nil
}

func productToFoodItem(p offProduct) *FoodItem {
	return &FoodItem{
		Barcode:     p.Code,
		Name:        p.ProductName,
		Brand:       p.Brands,
		ImageURL:    p.ImageURL,
		Calories:    p.Nutriments.EnergyKcal100g,
		ProteinG:    p.Nutriments.Proteins100g,
		CarbsG:      p.Nutriments.Carbs100g,
		FatG:        p.Nutriments.Fat100g,
		FiberG:      p.Nutriments.Fiber100g,
		ServingSize: p.ServingSize,
	}
}
