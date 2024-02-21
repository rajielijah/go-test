package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Product represents a product with ID, Name, MerchantID, and DateAdded fields
type Product struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	MerchantID string    `json:"merchantId"`
	DateAdded  time.Time `json:"dateAdded"`
}

// ProductsMap is a thread-safe map to store products
var (
	productsMap sync.Map
)

// Display all products for a given merchant
func displayProducts(w http.ResponseWriter, r *http.Request) {
	merchantID := r.URL.Query().Get("merchantId")
	var products []Product

	productsMap.Range(func(key, value interface{}) bool {
		product := value.(Product)
		if product.MerchantID == merchantID {
			products = append(products, product)
		}
		return true // continue iteration
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Create a new product
func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product.DateAdded = time.Now()
	productsMap.Store(product.ID, product)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// Edit an existing product
func editProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := productsMap.Load(product.ID); !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	product.DateAdded = time.Now() // Update the date added to reflect the edit time
	productsMap.Store(product.ID, product)

	json.NewEncoder(w).Encode(product)
}

// Delete an existing product
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if _, ok := productsMap.Load(id); !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	productsMap.Delete(id)
	fmt.Fprintf(w, "Deleted product with ID %s", id)
}

func main() {
	http.HandleFunc("/products", displayProducts)     // GET request with merchantId query param
	http.HandleFunc("/product/create", createProduct) // POST request
	http.HandleFunc("/product/edit", editProduct)     // PUT request
	http.HandleFunc("/product/delete", deleteProduct) // DELETE request with id query param

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
