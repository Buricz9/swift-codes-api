package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"swift-codes-api/internal/service"
)

type SwiftHandler struct {
	service service.SwiftService
}

func NewSwiftHandler(service service.SwiftService) *SwiftHandler {
	return &SwiftHandler{service: service}
}

func (h *SwiftHandler) GetSwiftCode(w http.ResponseWriter, r *http.Request) {
	swiftCodeParam := chi.URLParam(r, "swiftCode")

	result, err := h.service.GetSwiftCodeWithBranches(swiftCodeParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *SwiftHandler) GetSwiftCodesByCountry(w http.ResponseWriter, r *http.Request) {
	countryISO2 := chi.URLParam(r, "countryISO2")

	result, err := h.service.GetSwiftCodesByCountry(countryISO2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *SwiftHandler) CreateSwiftCode(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SwiftCode            string  `json:"swiftCode"`
		BankName             string  `json:"bankName"`
		Address              string  `json:"address"`
		CountryISO2          string  `json:"countryISO2"`
		CountryName          string  `json:"countryName"`
		IsHeadquarter        bool    `json:"isHeadquarter"`
		HeadquarterSwiftCode *string `json:"headquarterSwiftCode"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.CreateSwiftCode(service.CreateSwiftCodeInput{
		SwiftCode:            input.SwiftCode,
		BankName:             input.BankName,
		Address:              input.Address,
		CountryISO2:          input.CountryISO2,
		CountryName:          input.CountryName,
		IsHeadquarter:        input.IsHeadquarter,
		HeadquarterSwiftCode: input.HeadquarterSwiftCode,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Swift Code created successfully"}`))
}

func (h *SwiftHandler) DeleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	swiftCodeParam := chi.URLParam(r, "swiftCode")

	err := h.service.DeleteSwiftCode(swiftCodeParam)
	if err != nil {
		if err.Error() == fmt.Sprintf("swift code not found: %s", swiftCodeParam) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Swift Code deleted successfully"}`))
}
