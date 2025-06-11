//go:generate mockgen -source=municipality.go -destination=../../../tests/mock/domain/municipality.mock.go
package domain

import (
	"context"

	"g_gen/internal/domain/model"
)

type Municipality interface {
	FindAll(ctx context.Context) ([]*model.Municipality, error)
	FindByID(ctx context.Context, id int) (*model.Municipality, error)
	FindByPrefectureCode(ctx context.Context, prefectureCode string) ([]*model.Municipality, error)
}
