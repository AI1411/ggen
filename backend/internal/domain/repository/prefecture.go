//go:generate mockgen -source=prefecture.go -destination=../../../tests/mock/domain/prefecture.mock.go
package domain

import (
	"context"

	"g_gen/internal/domain/model"
)

type PrefectureRepository interface {
	FindAll(ctx context.Context) ([]*model.Prefecture, error)
	FindByCode(ctx context.Context, code string) (*model.Prefecture, error)
}
