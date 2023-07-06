package repository

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/yoratyo/go-redpanda/storage"
)

type Repository struct {
	st storage.Client
}

func NewRepository(st storage.Client) Repository {
	return Repository{
		st: st,
	}
}

func (r *Repository) Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.RunWith(r.st.DB())
}

func (r *Repository) InsertCrypto(ctx context.Context, c *Crypto) error {
	insertBuilder := r.Builder().
		Insert(c.GetTableName()).
		SetMap(c.GetInsertMap())

	_, err := insertBuilder.ExecContext(ctx)

	return err
}

func (r *Repository) UpdateCrypto(ctx context.Context, c *Crypto) error {
	_, err := r.Builder().
		Update(c.GetTableName()).
		Where(sq.Eq{"code": c.GetCode()}).
		SetMap(c.GetUpdateMap()).
		ExecContext(ctx)

	return err
}

func (r *Repository) GetCryptoById(ctx context.Context, code string) (*Crypto, error) {
	crypto := &Crypto{}
	query := r.Builder().Select("*").From(crypto.GetTableName()).Where(sq.Eq{"code": code})

	entities, err := r.getByBuilder(ctx, &query)
	if err != nil {
		return nil, err
	}
	if len(entities) == 1 {
		crypto = entities[0]
	} else {
		crypto = nil
	}

	return crypto, nil
}

func (r *Repository) List(ctx context.Context, page uint64, pageSize uint64, filter CryptoListFilter) ([]*Crypto, error) {
	entity := Crypto{}
	builder := r.Builder().
		Select("*").
		From(entity.GetTableName()).
		Limit(pageSize).
		Offset(page * pageSize)

	builder = r.applyListFilter(builder, filter)

	return r.getByBuilder(ctx, &builder)
}

func (r *Repository) ListSize(ctx context.Context, filter CryptoListFilter) (size int, err error) {
	entity := Crypto{}

	builder := r.Builder().
		Select("count(code)").
		From(entity.GetTableName())

	builder = r.applyListFilter(builder, filter)

	row := builder.QueryRowContext(ctx)

	err = row.Scan(&size)
	if err != nil {
		return size, fmt.Errorf("scan err: %s", err)
	}

	return
}

func (r *Repository) applyListFilter(builder sq.SelectBuilder, filter CryptoListFilter) sq.SelectBuilder {
	if filter.Code != "" {
		builder = builder.Where(sq.Eq{"code": filter.Code})
	} else {
		if filter.Name != "" {
			builder = builder.Where(sq.Eq{"name": filter.Name})
		}

		if filter.Category != "" {
			builder = builder.Where(sq.Eq{"category": filter.Category})
		}

		if filter.Algorithm != "" {
			builder = builder.Where(sq.Eq{"algorithm": filter.Algorithm})
		}

		if filter.Platform != "" {
			builder = builder.Where(sq.Eq{"platform": filter.Platform})
		}

		if filter.Industry != "" {
			builder = builder.Where(sq.Eq{"industry": filter.Industry})
		}

		if filter.Types != "" {
			builder = builder.Where(sq.Eq{"types": filter.Types})
		}

		if filter.Mineable != nil {
			builder = builder.Where(sq.Eq{"mineable": filter.Mineable})
		}

		if filter.Audited != nil {
			builder = builder.Where(sq.Eq{"audited": filter.Audited})
		}

		if filter.PriceMin != 0 {
			builder = builder.Where(sq.GtOrEq{"price": filter.PriceMin})
		}

		if filter.PriceMax != 0 {
			builder = builder.Where(sq.LtOrEq{"price": filter.PriceMax})
		}
	}

	return builder
}

func (r *Repository) getByBuilder(ctx context.Context, builder *sq.SelectBuilder) (result []*Crypto, err error) {
	rows, err := builder.QueryContext(ctx)

	switch {
	case err == sql.ErrNoRows:
		return result, nil
	case err != nil:
		return
	}

	defer func(rows *sql.Rows) { _ = rows.Close() }(rows)

	entity := Crypto{}

	for rows.Next() {
		current := entity
		err = rows.Scan(current.GetSelectColumnPointers()...)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %s", err)
		}
		result = append(result, &current)
	}

	err = rows.Err()

	return
}
