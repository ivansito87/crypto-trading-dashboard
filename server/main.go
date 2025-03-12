package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// PostgreSQL Connection String (Update credentials)
const dbConnStr = "postgres://postgres:password@localhost:5432/crypto_trading?sslmode=disable"

var db *sql.DB

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Mock crypto price data
var cryptoPrices = map[string]float64{
	"BTC": 50000.00,
	"ETH": 3000.00,
	"ADA": 1.50,
}

var broadcast = make(chan map[string]float64)

// WebSocket Handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	for {
		prices := <-broadcast
		err := conn.WriteJSON(prices)
		if err != nil {
			log.Println("WebSocket error:", err)
			break
		}
	}
}

// Update Prices in Real Time
func updatePrices() {
	for {
		for symbol := range cryptoPrices {
			cryptoPrices[symbol] += rand.Float64()*200 - 100 // Random price change
		}
		broadcast <- cryptoPrices
		time.Sleep(2 * time.Second)
	}
}

// Order Struct
type Order struct {
	ID       int     `json:"id"`
	Symbol   string  `json:"symbol"`
	Type     string  `json:"type"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	DateTime string  `json:"datetime"`
}

// Handles Buy/Sell Orders & Stores in PostgreSQL
func handleOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate symbol
	if _, exists := cryptoPrices[order.Symbol]; !exists {
		http.Error(w, "Invalid cryptocurrency symbol", http.StatusBadRequest)
		return
	}

	// Assign order price & timestamp
	order.Price = cryptoPrices[order.Symbol]
	order.DateTime = time.Now().Format(time.RFC3339)

	// Insert into PostgreSQL
	query := `INSERT INTO trades (symbol, type, amount, price) VALUES ($1, $2, $3, $4) RETURNING id, datetime`
	err = db.QueryRow(query, order.Symbol, order.Type, order.Amount, order.Price).Scan(&order.ID, &order.DateTime)
	if err != nil {
		http.Error(w, "Failed to save trade", http.StatusInternalServerError)
		log.Println("DB Insert Error:", err)
		return
	}

	log.Printf("Trade Executed: %+v\n", order)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// Fetch Trade History from PostgreSQL
func getTradeHistory(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, symbol, type, amount, price, datetime FROM trades ORDER BY datetime DESC")
	if err != nil {
		http.Error(w, "Failed to fetch trade history", http.StatusInternalServerError)
		log.Println("DB Query Error:", err)
		return
	}
	defer rows.Close()

	var trades []Order
	for rows.Next() {
		var trade Order
		err := rows.Scan(&trade.ID, &trade.Symbol, &trade.Type, &trade.Amount, &trade.Price, &trade.DateTime)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		trades = append(trades, trade)
	}

	if len(trades) == 0 || trades == nil {
		trades = []Order{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trades)
}

func main() {
	// Initialize database connection
	var err error
	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Check DB connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	// end of db connection

	// Start WebSocket price updates with go routine
	go updatePrices()

	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/ws", handleWebSocket)
	router.HandleFunc("/order", handleOrder).Methods("POST")
	router.HandleFunc("/trades", getTradeHistory).Methods("GET")

	// Configure CORS options
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Allows all origins, in production specify origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	port := ":8080"
	fmt.Println("Crypto Trading Server running on", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
