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

func NewSwiftHandler(svc service.SwiftService) *SwiftHandler {
	return &SwiftHandler{service: svc}
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
