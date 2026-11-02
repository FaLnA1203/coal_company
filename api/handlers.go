package api

import (
	"coal_company/api/dto"
	"coal_company/company"
	"coal_company/company/miners"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
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

func (handlers *HTTPHandler) BuyMinerHandler(w http.ResponseWriter, r *http.Request) {
	var m dto.InputDTOMiner
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		errDTO := ErrorDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	miner, err := handlers.company.HireMiner(m.MinerType)
	if err != nil {
		errDTO := ErrorDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		if errors.Is(err, company.ErrInsufficientFunds) {
			http.Error(w, errDTO.ToString(), http.StatusForbidden)
		} else if errors.Is(err, company.ErrUnknokwnMinerType) {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(miner.Info(), "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

func (handlers *HTTPHandler) BuyEquipmentHandler(w http.ResponseWriter, r *http.Request) {
	var eq dto.InputDTOEquipment
	if err := json.NewDecoder(r.Body).Decode(&eq); err != nil {
		errDTO := ErrorDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	_, err := handlers.company.BuyEquipment(eq.EquipmentType)
	if err != nil {
		errDTO := ErrorDTO{
			Error: err.Error(),
			Time:  time.Now(),
		}
		if errors.Is(err, company.ErrInsufficientFunds) {
			http.Error(w, errDTO.ToString(), http.StatusForbidden)
		} else if errors.Is(err, company.ErrUnknokwnEquipmentType) {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	eqDTO := dto.NewEquipmentDTO(handlers.company.GetEquipment())
	b, err := json.MarshalIndent(eqDTO, "", "    ")
	if err != nil {
		panic(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

func (handlers *HTTPHandler) HandlerGetMiners(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("miner")
	lm := make(map[miners.MinersType]map[uuid.UUID]miners.Miner)
	switch q {
	case "Little":
		lm[miners.MinersType(q)] = handlers.company.GetMinersByType(miners.MinersType(q))
	case "Middle":
		lm[miners.MinersType(q)] = handlers.company.GetMinersByType(miners.MinersType(q))
	case "Strong":
		lm[miners.MinersType(q)] = handlers.company.GetMinersByType(miners.MinersType(q))
	default:
		lm = handlers.company.GetAllMiners()
	}

	b, err := json.MarshalIndent(lm, "", "    ")
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

func (handlers *HTTPHandler) HandlerGetEquipment(w http.ResponseWriter, r *http.Request) {
	equipment := handlers.company.GetEquipment()
	dtoequipment := dto.NewEquipmentDTO(equipment)
	b, err := json.MarshalIndent(dtoequipment, "", "  ")
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

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
