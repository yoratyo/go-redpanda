package service

import (
	"context"
	"log"
	"math"

	"github.com/yoratyo/go-redpanda/kafka"
	"github.com/yoratyo/go-redpanda/repository"
)

type Service struct {
	repo      repository.Repository
	publisher *kafka.Publisher
}

func NewService(
	repo repository.Repository,
	publisher *kafka.Publisher,
) Service {
	return Service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) PublishCrypto(ctx context.Context, c CryptoDTO) error {
	return s.publisher.PostCrypto(ctx, *c.ToEventModel())
}

func (s *Service) ConsumeCrypto(ctx context.Context, value []byte) error {
	var msg kafka.CryptoModel
	err := msg.DeSerialize(value)
	if err != nil {
		return err
	}
	log.Printf("Incoming event: %+v\n", msg)

	return s.UpsertCrypto(ctx, CryptoDTO{
		Code:      msg.Code,
		Name:      msg.Name,
		Category:  msg.Category,
		Algorithm: msg.Algorithm,
		Platform:  msg.Platform,
		Industry:  msg.Industry,
		Types:     msg.Types,
		Mineable:  msg.Mineable,
		Audited:   msg.Audited,
		Price:     msg.Price,
	})
}

func (s *Service) UpsertCrypto(ctx context.Context, c CryptoDTO) error {
	crypto, err := s.repo.GetCryptoById(ctx, c.Code)
	if err != nil {
		return err
	}
	if crypto == nil {
		return s.repo.InsertCrypto(ctx, c.ToDAO())
	}

	return s.repo.UpdateCrypto(ctx, c.ToDAO())
}

func (s *Service) ListCrypto(ctx context.Context, req ListCryptoRequest) (*ListCryptoResponse, error) {
	var page uint64 = 1
	if req.Page > 0 {
		page = uint64(req.Page)
	}

	var pageSize uint64 = 10
	if req.PageSize > 0 {
		pageSize = uint64(req.PageSize)
	}

	filter := repository.CryptoListFilter{
		Code:      req.Code,
		Name:      req.Name,
		Category:  req.Category,
		Algorithm: req.Algorithm,
		Platform:  req.Platform,
		Industry:  req.Industry,
		Types:     req.Types,
		Mineable:  req.Mineable,
		Audited:   req.Audited,
		PriceMin:  req.PriceMin,
		PriceMax:  req.PriceMax,
	}

	list, err := s.repo.List(ctx, (page - 1), pageSize, filter)
	if err != nil {
		return nil, err
	}

	listSize, err := s.repo.ListSize(ctx, filter)
	if err != nil {
		return nil, err
	}

	cryptos := make([]CryptoDTO, len(list))
	for i, c := range list {
		cryptos[i] = CryptoDTO{
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

	return &ListCryptoResponse{
		Page:       uint32(page),
		PageSize:   uint32(pageSize),
		TotalPages: uint32(math.Ceil(float64(listSize) / float64(pageSize))),
		TotalItems: uint32(listSize),
		Cryptos:    cryptos,
	}, nil
}
