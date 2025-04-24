package handler

import (
	"encoding/json"
	"errors"
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
	ctx := r.Context()
	code := chi.URLParam(r, "swiftCode")

	result, err := h.service.GetSwiftCodeWithBranches(ctx, code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, err.Error())
		} else {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (h *SwiftHandler) GetSwiftCodesByCountry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	country := chi.URLParam(r, "countryISO2")

	result, err := h.service.GetSwiftCodesByCountry(ctx, country)
	if err != nil {
		if errors.Is(err, service.ErrNoCountryCodes) {
			writeJSONError(w, http.StatusNotFound, err.Error())
		} else {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (h *SwiftHandler) CreateSwiftCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input struct {
		SwiftCode            string  `json:"swiftCode"`
		BankName             string  `json:"bankName"`
		Address              string  `json:"address"`
		CountryISO2          string  `json:"countryISO2"`
		CountryName          string  `json:"countryName"`
		IsHeadquarter        bool    `json:"isHeadquarter"`
		HeadquarterSwiftCode *string `json:"headquarterSwiftCode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.service.CreateSwiftCode(ctx, service.CreateSwiftCodeInput{
		SwiftCode:            input.SwiftCode,
		BankName:             input.BankName,
		Address:              input.Address,
		CountryISO2:          input.CountryISO2,
		CountryName:          input.CountryName,
		IsHeadquarter:        input.IsHeadquarter,
		HeadquarterSwiftCode: input.HeadquarterSwiftCode,
	})
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Swift Code created successfully"})
}

func (h *SwiftHandler) DeleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	code := chi.URLParam(r, "swiftCode")

	err := h.service.DeleteSwiftCode(ctx, code)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, err.Error())
		} else {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Swift Code deleted successfully"})
}
