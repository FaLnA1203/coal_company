package api

import (
	"coal_company/api/dto"
	"coal_company/company"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPHandler struct {
	company     *company.Company
	closeServer func() error
}

func NewAPIHandler(company *company.Company) HTTPHandler {
	return HTTPHandler{
		company: company,
	}
}

func (handlers *HTTPHandler) SetCloseServerFunc(f func() error) {
	handlers.closeServer = f
}

func (handlers *HTTPHandler) PriceMinersHandler(w http.ResponseWriter, r *http.Request) {
	var price *dto.DTOPriceMiner
	price = dto.NewDTOPriceMiner()
	b, err := json.MarshalIndent(price, "", "  ")
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

func (handlers *HTTPHandler) PriceEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	var price *dto.EquipmentPricesDTO
	price = dto.NewEquipmentPricesDTO()
	b, err := json.MarshalIndent(price, "", "  ")
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

func (handlers *HTTPHandler) BuyMinerHandler(w http.ResponseWriter, r *http.Request) {}

func (handlers *HTTPHandler) BuyEquipmentHandler(w http.ResponseWriter, r *http.Request) {}

func (handlers *HTTPHandler) HandlerGetMiners(w http.ResponseWriter, r *http.Request) {}

func (handlers *HTTPHandler) HandlerGetEquipment(w http.ResponseWriter, r *http.Request) {}

func (handlers *HTTPHandler) HandlerCloseCompany(w http.ResponseWriter, r *http.Request) {
	stat, err := handlers.company.Complete()
	if err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, company.ErrNotAllEquipmentPurchased) {
			http.Error(w, errDTO.ToString(), http.StatusForbidden)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	statOutputDTO := dto.NewCompanyStatisticsDTO(stat)
	b, err := json.MarshalIndent(statOutputDTO, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

func (handlers *HTTPHandler) HandlerInfoCompany(w http.ResponseWriter, r *http.Request) {
	var stat *dto.CompanyStatisticsDTO
	stat = dto.NewCompanyStatisticsDTO(handlers.company.GetStatistics())
	b, err := json.MarshalIndent(stat, "", "  ")
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}
