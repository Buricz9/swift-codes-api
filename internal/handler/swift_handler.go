package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
	"swift-codes-api/internal/service"
)

type CreateSwiftRequest struct {
	SwiftCode            string  `json:"swiftCode"`
	BankName             string  `json:"bankName"`
	Address              string  `json:"address"`
	CountryISO2          string  `json:"countryISO2"`
	CountryName          string  `json:"countryName"`
	IsHeadquarter        bool    `json:"isHeadquarter"`
	HeadquarterSwiftCode *string `json:"headquarterSwiftCode"`
}

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

	var req CreateSwiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONMessage(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if errs := req.Validate(); len(errs) > 0 {
		parts := make([]string, 0, len(errs))
		for field, msg := range errs {
			parts = append(parts, field+": "+msg)
		}
		writeJSONMessage(w, http.StatusBadRequest, "validation failed: "+strings.Join(parts, "; "))
		return
	}

	input := service.CreateSwiftCodeInput{
		SwiftCode:            req.SwiftCode,
		BankName:             req.BankName,
		Address:              req.Address,
		CountryISO2:          req.CountryISO2,
		CountryName:          req.CountryName,
		IsHeadquarter:        req.IsHeadquarter,
		HeadquarterSwiftCode: req.HeadquarterSwiftCode,
	}

	if err := h.service.CreateSwiftCode(ctx, input); err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			writeJSONMessage(w, http.StatusConflict, err.Error())
		} else {
			writeJSONMessage(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSONMessage(w, http.StatusCreated, "Swift Code created successfully")
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
	_ = json.NewEncoder(w).Encode(MessageResponse{
		Message: "Swift Code deleted successfully",
	})
}

func (r *CreateSwiftRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if len(r.SwiftCode) != 8 && len(r.SwiftCode) != 11 {
		errs["swiftCode"] = "must be 8 or 11 chars"
	}
	if r.BankName == "" {
		errs["bankName"] = "cannot be empty"
	}
	if r.Address == "" {
		errs["address"] = "cannot be empty"
	}
	if len(r.CountryISO2) != 2 || strings.ToUpper(r.CountryISO2) != r.CountryISO2 {
		errs["countryISO2"] = "must be 2 uppercase letters"
	}
	if r.CountryName == "" {
		errs["countryName"] = "cannot be empty"
	}
	if r.HeadquarterSwiftCode != nil {
		hq := *r.HeadquarterSwiftCode
		if len(hq) != 11 {
			errs["headquarterSwiftCode"] = "if provided, must be 11 chars"
		}
	}
	return errs
}
