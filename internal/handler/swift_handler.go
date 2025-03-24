package handler

import (
	"encoding/json"
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

	result, err := h.service.GetSwiftCode(swiftCodeParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var headquarterSwiftCode *string
	if result.HeadquarterSwiftCode.Valid {
		headquarterSwiftCode = &result.HeadquarterSwiftCode.String
	}

	response := map[string]interface{}{
		"id":                   result.ID,
		"swiftCode":            result.SwiftCode,
		"bankName":             result.BankName,
		"address":              result.Address,
		"countryISO2":          result.CountryISO2,
		"countryName":          result.CountryName,
		"isHeadquarter":        result.IsHeadquarter,
		"headquarterSwiftCode": headquarterSwiftCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *SwiftHandler) GetSwiftCodesByCountry(w http.ResponseWriter, r *http.Request) {
	countryISO2 := chi.URLParam(r, "countryISO2")

	swiftCodes, err := h.service.GetSwiftCodesByCountry(countryISO2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var response []map[string]interface{}
	for _, sc := range swiftCodes {
		var headquarterSwiftCode *string
		if sc.HeadquarterSwiftCode.Valid {
			headquarterSwiftCode = &sc.HeadquarterSwiftCode.String
		}

		item := map[string]interface{}{
			"id":                   sc.ID,
			"swiftCode":            sc.SwiftCode,
			"bankName":             sc.BankName,
			"address":              sc.Address,
			"countryISO2":          sc.CountryISO2,
			"countryName":          sc.CountryName,
			"isHeadquarter":        sc.IsHeadquarter,
			"headquarterSwiftCode": headquarterSwiftCode,
		}

		response = append(response, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
