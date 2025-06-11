package datastore

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"g_gen/internal/domain/model"
	"g_gen/internal/domain/query"
	domain "g_gen/internal/domain/repository"
	myerrors "g_gen/internal/errors"
	"g_gen/internal/infra/db"
)

type prefectureRepository struct {
	client db.Client
	query  *query.Query
}

func NewPrefectureRepository(
	ctx context.Context,
	client db.Client,
) domain.PrefectureRepository {
	return &prefectureRepository{
		client: client,
		query:  query.Use(client.Conn(ctx)),
	}
}

func (r *prefectureRepository) FindAll(ctx context.Context) ([]*model.Prefecture, error) {
	prefectures, err := r.query.WithContext(ctx).Prefecture.Find()
	if err != nil {
		return nil, err
	}

	return prefectures, nil
}

func (r *prefectureRepository) FindByCode(ctx context.Context, code string) (*model.Prefecture, error) {
	prefecture, err := r.query.WithContext(ctx).
		Prefecture.
		Where(r.query.Prefecture.Code.Eq(code)).
		Preload(r.query.Prefecture.Municipalities.On(r.query.Municipality.IsActive.Is(true))).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &myerrors.APIError{
				Code:    myerrors.PrefectureNotFoundError,
				Message: myerrors.PrefectureNotFoundErrorMessage,
			}
		}

		return nil, err
	}

	return prefecture, nil
}
