//go:generate mockgen -source=prefecture_usecase.go -destination=../../tests/mock/usecase/prefecture_usecase.mock.go
package usecase

import (
	"context"

	"g_gen/internal/domain/model"
	domain "g_gen/internal/domain/repository"
)

type PrefectureUseCase interface {
	ListPrefectures(ctx context.Context) ([]*model.Prefecture, error)
	GetPrefectureByCode(ctx context.Context, code string) (*model.Prefecture, error)
}

type prefectureUseCase struct {
	prefectureRepository domain.PrefectureRepository
}

func NewPrefectureUseCase(
	prefectureRepository domain.PrefectureRepository,
) PrefectureUseCase {
	return &prefectureUseCase{
		prefectureRepository: prefectureRepository,
	}
}

func (u *prefectureUseCase) ListPrefectures(ctx context.Context) ([]*model.Prefecture, error) {
	prefectures, err := u.prefectureRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return prefectures, nil
}

func (u *prefectureUseCase) GetPrefectureByCode(ctx context.Context, code string) (*model.Prefecture, error) {
	prefecture, err := u.prefectureRepository.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return prefecture, nil
}
