package service

import (
	"github.com/yoratyo/go-redpanda/kafka"
	"github.com/yoratyo/go-redpanda/repository"
)

type CryptoDTO struct {
	Code      string  `json:"code" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Category  string  `json:"category"`
	Algorithm string  `json:"algorithm"`
	Platform  string  `json:"platform"`
	Industry  string  `json:"industry"`
	Types     string  `json:"types" validate:"required"`
	Mineable  bool    `json:"mineable" validate:"required"`
	Audited   bool    `json:"audited" validate:"required"`
	Price     float64 `json:"price" validate:"required"`
}

func (c *CryptoDTO) ToDAO() *repository.Crypto {
	return &repository.Crypto{
		Code:      c.Code,
		Name:      c.Name,
		Category:  c.Category,
		Algorithm: c.Algorithm,
		Platform:  c.Platform,
		Industry:  c.Industry,
		Types:     c.Types,
		Mineable:  c.Mineable,
		Audited:   c.Audited,
		Price:     c.Price,
	}
}

func (c *CryptoDTO) ToEventModel() *kafka.CryptoModel {
	return &kafka.CryptoModel{
		Code:      c.Code,
		Name:      c.Name,
		Category:  c.Category,
		Algorithm: c.Algorithm,
		Platform:  c.Platform,
		Industry:  c.Industry,
		Types:     c.Types,
		Mineable:  c.Mineable,
		Audited:   c.Audited,
		Price:     c.Price,
	}
}

type ListCryptoRequest struct {
	Page      uint32  `schema:"page"`
	PageSize  uint32  `schema:"pageSize"`
	Code      string  `schema:"code"`
	Name      string  `schema:"name"`
	Category  string  `schema:"category"`
	Algorithm string  `schema:"algorithm"`
	Platform  string  `schema:"platform"`
	Industry  string  `schema:"industry"`
	Types     string  `schema:"types"`
	Mineable  *bool   `schema:"mineable"`
	Audited   *bool   `schema:"audited"`
	PriceMin  float64 `schema:"priceMin"`
	PriceMax  float64 `schema:"priceMax"`
}

type ListCryptoResponse struct {
	Page       uint32      `json:"page"`
	PageSize   uint32      `json:"pageSize"`
	TotalPages uint32      `json:"totalPages"`
	TotalItems uint32      `json:"totalItems"`
	Cryptos    []CryptoDTO `json:"cryptos"`
}
