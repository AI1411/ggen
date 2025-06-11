package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"g_gen/internal/infra/logger"
	"g_gen/internal/usecase"
)

type PrefectureHandler interface {
	ListPrefectures(c *gin.Context)
	GetPrefecture(c *gin.Context)
}

type prefectureHandler struct {
	appLogger         *logger.Logger
	prefectureUseCase usecase.PrefectureUseCase
}

func NewPrefectureHandler(
	l *logger.Logger,
	prefectureUseCase usecase.PrefectureUseCase,
) PrefectureHandler {
	return &prefectureHandler{
		appLogger:         l,
		prefectureUseCase: prefectureUseCase,
	}
}

type PrefectureResponse struct {
	ID   int32  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// ListPrefectures @title 都道府県一覧取得
// @id ListPrefectures
// @tags prefectures
// @accept json
// @produce json
// @Summary 都道府県一覧取得
// @Success 200 {array} PrefectureResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Security BearerAuth
// @Description 都道府県の一覧を取得します。
// @Router /prefectures [get]
func (h *prefectureHandler) ListPrefectures(c *gin.Context) {
	prefectures, err := h.prefectureUseCase.ListPrefectures(c.Request.Context())
	if err != nil {
		handleError(c, err, h.appLogger, "Failed to list prefectures")

		return
	}

	var response []*PrefectureResponse
	for _, p := range prefectures {
		response = append(response, &PrefectureResponse{
			Code: p.Code,
			Name: p.Name,
		})
	}

	c.JSON(http.StatusOK, response)
}

type GetPrefectureResponse struct {
	ID             int32           `json:"id"`
	Name           string          `json:"name"`
	Municipalities []*Municipality `json:"municipalities"`
}

type Municipality struct {
	ID                    int32  `json:"id"`
	PrefectureCode        string `json:"prefecture_code"`
	OrganizationCode      string `json:"organization_code"`
	PrefectureNameKanji   string `json:"prefecture_name_kanji"`
	MunicipalityNameKanji string `json:"municipality_name_kanji"`
	PrefectureNameKana    string `json:"prefecture_name_kana"`
	MunicipalityNameKana  string `json:"municipality_name_kana"`
	IsActive              bool   `json:"is_active"`
}

type GetPrefectureRequest struct {
	Code string `form:"code" binding:"required,numeric"`
}

// GetPrefecture @title 都道府県詳細取得
// @id GetPrefecture
// @tags prefectures
// @accept json
// @produce json
// @Param code path string true "都道府県コード"
// @Description 都道府県コードを指定して、都道府県の詳細情報を取得します。
// @Summary 都道府県詳細取得
// @Success 200 {object} PrefectureResponse
// @Failure 400 {object} ErrorResponseDetail
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /prefectures/{code} [get]
func (h *prefectureHandler) GetPrefecture(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")
	var req GetPrefectureRequest
	req.Code = code
	if err := c.ShouldBindUri(&req); err != nil {
		handleValidationError(c, err, h.appLogger, "invalid prefecture code")

		return
	}

	prefecture, err := h.prefectureUseCase.GetPrefectureByCode(ctx, code)
	if err != nil {
		handleError(c, err, h.appLogger, "failed to get prefecture")

		return
	}

	response := &GetPrefectureResponse{
		ID:             prefecture.ID,
		Name:           prefecture.Name,
		Municipalities: make([]*Municipality, len(prefecture.Municipalities)),
	}

	for i, municipality := range prefecture.Municipalities {
		response.Municipalities[i] = &Municipality{
			ID:                    municipality.ID,
			PrefectureCode:        municipality.PrefectureCode,
			OrganizationCode:      municipality.OrganizationCode,
			PrefectureNameKanji:   municipality.PrefectureNameKanji,
			MunicipalityNameKanji: municipality.MunicipalityNameKanji,
			PrefectureNameKana:    municipality.PrefectureNameKana,
			MunicipalityNameKana:  municipality.MunicipalityNameKana,
			IsActive:              municipality.IsActive,
		}
	}

	c.JSON(http.StatusOK, response)
}
