package routes

import (
	"encoding/json"
	"example-api-chi/internal/db"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CustomerRoutes(r chi.Router) {
	// Probably want to do some middleware to create a more consistent serializer but this should work for now
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/", HandleGetCustomers)    // GET /customer
	r.Post("/", HandleCreateCustomer) // POST /customer

	r.Get("/{id}", HandleGetCustomerByID) // GET /customer/{id}
	// TODO:
	// PATCH /customer/{id}
	// DELETE /customer/{id}
}

func HandleGetCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := db.GetAllCustomers()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

func HandleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	customerID := chi.URLParam(r, "id")
	customer, err := db.GetCustomerByID(customerID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req db.InsertRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	response, err := db.CreateCustomerFullResponse(db.InsertRequest{Name: req.Name})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
