package repository

import (
	"encoding/json"
	"reflect"
	"time"
)

const (
	table = "cryptocurrency"
)

type Crypto struct {
	Code      string
	Name      string
	Category  string
	Algorithm string
	Platform  string
	Industry  string
	Types     string
	Mineable  bool
	Audited   bool
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CryptoListFilter struct {
	Code      string
	Name      string
	Category  string
	Algorithm string
	Platform  string
	Industry  string
	Types     string
	Mineable  *bool
	Audited   *bool
	PriceMin  float64
	PriceMax  float64
}

func (c *Crypto) GetCode() string      { return c.Code }
func (c *Crypto) GetTableName() string { return table }

func (c *Crypto) GetInsertMap() map[string]any {
	b, _ := json.Marshal(&c)
	var result map[string]any
	_ = json.Unmarshal(b, &result)
	delete(result, "CreatedAt")
	delete(result, "UpdatedAt")

	return result
}

func (c *Crypto) GetUpdateMap() map[string]any {
	b, _ := json.Marshal(&c)
	var result map[string]any
	_ = json.Unmarshal(b, &result)
	delete(result, "Code")
	delete(result, "CreatedAt")
	delete(result, "UpdatedAt")

	return result
}

func (s *Crypto) GetSelectColumnPointers() []any {
	e := reflect.ValueOf(s).Elem()
	result := make([]any, e.NumField())
	for i := 0; i < e.NumField(); i++ {
		result[i] = e.Field(i).Addr().Interface()
	}
	return result
}
