package routes

import (
	"encoding/json"
	"example-api-chi/internal/model"
	svc "example-api-chi/internal/service"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type CustomerHandler struct {
	svc *svc.CustomerService
}

func NewCustomerHandler(service *svc.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		svc: service,
	}
}

func CustomerRoutes(r chi.Router, service *svc.CustomerService) {
	handler := NewCustomerHandler(service)
	handler.Routes(r)
}

func (h *CustomerHandler) Routes(r chi.Router) {
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Get("/", h.HandleGetCustomers)    // GET /customer
	r.Post("/", h.HandleCreateCustomer) // POST /customer

	r.Get("/{id}", h.HandleGetCustomerById)       // GET /customer/{id}
	r.Patch("/{id}", h.HandleUpdateCustomerById)  // PATCH /customer/{id}
	r.Delete("/{id}", h.HandleDeleteCustomerById) // DELETE /customer/{id}
}

func (h *CustomerHandler) HandleGetCustomers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customers, err := h.svc.GetAllCustomers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customers)
}

func (h *CustomerHandler) HandleGetCustomerById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerId := chi.URLParam(r, "id")
	customer, err := h.svc.GetCustomerByID(ctx, customerId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) HandleCreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req model.Customer

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Name == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body.  Name field cannot be blank"})
		return
	}

	response, err := h.svc.CreateCustomer(ctx, req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *CustomerHandler) HandleUpdateCustomerById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req model.Customer
	customerId := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if req.Name == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body.  Name field cannot be blank"})
		return
	}

	customer, err := h.svc.UpdateCustomer(ctx, customerId, req.Name)
	if err != nil {
		if strings.Contains(err.Error(), "cannot be empty") {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customer)
}

func (h *CustomerHandler) HandleDeleteCustomerById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	customerId := chi.URLParam(r, "id")

	err := h.svc.DeleteCustomer(ctx, customerId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusNoContent)
}
