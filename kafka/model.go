package kafka

import (
	"encoding/json"

	"github.com/yoratyo/go-redpanda/repository"
)

type CryptoModel struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Category  string  `json:"category"`
	Algorithm string  `json:"algorithm"`
	Platform  string  `json:"platform"`
	Industry  string  `json:"industry"`
	Types     string  `json:"types"`
	Mineable  bool    `json:"minable"`
	Audited   bool    `json:"audited"`
	Price     float64 `json:"price"`
}

func (m *CryptoModel) Serialize() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (m *CryptoModel) DeSerialize(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	return nil
}

func (c *CryptoModel) ToDAO() *repository.Crypto {
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
