package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Product struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}
type AddToCartRequest struct {
	UserId    int `json:"user_id"`
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
type CartItem struct {
	ProductId int    `json:"product_id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}
type RemoveFromCartRequest struct {
	UserId    int `json:"user_id"`
	ProductId int `json:"product_id"`
}
type CheckoutRequest struct {
	UserId int `json:"user_id"`
}

var productsCounter = 0
var products = make(map[int]Product)
var cart = make(map[int][]CartItem)

func addProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var p Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Status Bad Request", http.StatusBadRequest)
	}
	fmt.Println(p)
	p.Id = productsCounter
	products[productsCounter] = p
	productsCounter++
	fmt.Println(products)
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(products)
	json.NewEncoder(w).Encode(p)

}
func getProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "StatusMethodNotAllowed", http.StatusMethodNotAllowed)
		return
	}
	productList := []Product{}
	for _, val := range products {
		productList = append(productList, val)

	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(productList)
}
func addtocart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
	var req AddToCartRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Payload", http.StatusBadRequest)
		return
	}
	fmt.Println(req)
	productId := req.ProductId

	p, ok := products[productId]
	if !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return

	}
	cartItem := CartItem{ProductId: p.Id, Name: p.Name, Quantity: req.Quantity, Price: p.Price}
	cart[req.UserId] = append(cart[req.UserId], cartItem)

	fmt.Println(cart)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart[req.UserId])

}
func getCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Status Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	userId, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	fmt.Println(userId)
	val, ok := cart[userId]
	emptyResponse := []CartItem{}
	if !ok {
		json.NewEncoder(w).Encode(emptyResponse)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(val)

}
func removeFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Status Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	removeFromCartPayload := RemoveFromCartRequest{}

	json.NewDecoder(r.Body).Decode(&removeFromCartPayload)
	fmt.Println(removeFromCartPayload)
	userId := removeFromCartPayload.UserId
	productId := removeFromCartPayload.ProductId
	cartProducts, ok := cart[userId]
	if !ok {
		http.Error(w, "UserId not found", http.StatusNotFound)
		return
	}
	targetindex := -1
	for index, val := range cartProducts {
		if val.ProductId == productId {
			targetindex = index
			break

		}
	}
	if targetindex == -1 {
		http.Error(w, "Product Id not found", http.StatusNotFound)
		return
	}
	cartProducts = append(cartProducts[0:targetindex], cartProducts[targetindex+1:]...)
	cart[userId] = cartProducts
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart[userId])

}
func checkOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Status Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	checkoutRequestPayload := CheckoutRequest{}

	err := json.NewDecoder(r.Body).Decode(&checkoutRequestPayload)
	if err != nil {
		http.Error(w, "Status Bad Request", http.StatusBadRequest)
		return
	}
	userId := checkoutRequestPayload.UserId
	cartitems, ok := cart[userId]
	if !ok {
		http.Error(w, "Status Not Found", http.StatusNotFound)
		return
	}
	totalPrice := 0

	for _, item := range cartitems {
		totalPrice = totalPrice + item.Price*item.Quantity
	}
	cart[userId] = []CartItem{}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalPrice)

}
func getCartSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Status Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	userId, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Status Bad Request", http.StatusBadRequest)
		return
	}
	val, ok := cart[userId]
	if !ok {
		http.Error(w, "Status Not Found", http.StatusNotFound)
		return
	}
	price := 0
	totalItems := 0
	for _, item := range val {
		price = price + item.Price*item.Quantity
		totalItems += item.Quantity
	}
	response := map[string]int{
		"user_id": userId, "total_items": totalItems, "total_price": price,
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}
func searchProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Status Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	maxPricestr := r.URL.Query().Get("max_price")
	minPricestr := r.URL.Query().Get("min_price")
	name := r.URL.Query().Get("name")
	maxprice := int(^uint(0) >> 1)
	minprice := 0
	if maxPricestr != "" {
		maxprice, _ = strconv.Atoi(maxPricestr)
	}
	if minPricestr != "" {
		minprice, _ = strconv.Atoi(minPricestr)
	}
	fmt.Println(name, maxprice, minprice)
	http.Error(w, "Status Method Not Allowed", http.StatusMethodNotAllowed)
}
func updateProducts( w http.ResponseWriter, r *http.Request){
	
}
func main() {
	fmt.Println("hello")
	http.HandleFunc("/addproduct", addProduct)
	http.HandleFunc("/products", getProducts)
	http.HandleFunc("/addtocart", addtocart)
	http.HandleFunc("/cart", getCart)
	http.HandleFunc("/cart/remove", removeFromCart)
	http.HandleFunc("/cart/checkout", checkOutHandler)
	http.HandleFunc("/cart/summary", getCartSummary)
	http.HandleFunc("/search/products", searchProducts)
	http.HandleFunc("/products/update", updateProducts)

	log.Fatal(http.ListenAndServe(":3000", nil))

}
