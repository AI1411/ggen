package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"g_gen/internal/domain/model"
	"g_gen/internal/handler"
	"g_gen/internal/infra/logger"
	mockusecase "g_gen/tests/mock/usecase"
)

func TestPrefectureHandler_ListPrefectures(t *testing.T) {
	tests := []struct {
		name       string
		mockSetup  func(mockUseCase *mockusecase.MockPrefectureUseCase)
		wantStatus int
		wantBody   func() string
	}{
		{
			name: "Success",
			mockSetup: func(mockUseCase *mockusecase.MockPrefectureUseCase) {
				prefectures := expectedPrefectureListModel()
				mockUseCase.EXPECT().ListPrefectures(gomock.Any()).Return(prefectures, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func() string {
				response := expectedListPrefecturesResponse()
				responseJSON, _ := json.Marshal(response)
				return string(responseJSON)
			},
		},
		{
			name: "Error",
			mockSetup: func(mockUseCase *mockusecase.MockPrefectureUseCase) {
				mockUseCase.EXPECT().ListPrefectures(gomock.Any()).Return(nil, errors.New("database error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			appLogger := logger.New(logger.DefaultConfig())

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req := httptest.NewRequest(http.MethodGet, "/prefectures", nil)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = req

			uc := mockusecase.NewMockPrefectureUseCase(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(uc)
			}

			mockHandler := handler.NewPrefectureHandler(appLogger, uc)
			mockHandler.ListPrefectures(c)

			a.Equal(tt.wantStatus, rec.Code)

			if tt.wantBody != nil {
				wantBody := tt.wantBody()
				resBody := rec.Body.String()
				if !cmp.Equal(wantBody, resBody) {
					t.Errorf("diff: %s", cmp.Diff(wantBody, resBody))
				}
			}
		})
	}
}

func expectedPrefectureListModel() []*model.Prefecture {
	return []*model.Prefecture{
		{
			ID:   1,
			Name: "北海道",
			Code: "01",
			Municipalities: []model.Municipality{
				{
					ID:                    1,
					PrefectureCode:        "01",
					OrganizationCode:      "011002",
					PrefectureNameKanji:   "北海道",
					MunicipalityNameKanji: "札幌市",
					PrefectureNameKana:    "ﾎｯｶｲﾄﾞｳ",
					MunicipalityNameKana:  "ｻｯﾎﾟﾛｼ",
					IsActive:              true,
				},
			},
		},
		{
			ID:   13,
			Name: "東京都",
			Code: "13",
			Municipalities: []model.Municipality{
				{
					ID:                    2,
					PrefectureCode:        "13",
					OrganizationCode:      "131016",
					PrefectureNameKanji:   "東京都",
					MunicipalityNameKanji: "千代田区",
					PrefectureNameKana:    "ﾄｳｷｮｳﾄ",
					MunicipalityNameKana:  "ﾁﾖﾀﾞｸ",
					IsActive:              true,
				},
			},
		},
	}
}

func expectedListPrefecturesResponse() []*handler.PrefectureResponse {
	return []*handler.PrefectureResponse{
		{
			Code: "01",
			Name: "北海道",
		},
		{
			Code: "13",
			Name: "東京都",
		},
	}
}

func TestPrefectureHandler_GetPrefecture(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		mockSetup  func(mockUseCase *mockusecase.MockPrefectureUseCase)
		wantStatus int
		wantBody   func() string
	}{
		{
			name: "Success",
			code: "01",
			mockSetup: func(mockUseCase *mockusecase.MockPrefectureUseCase) {
				prefecture := expectedPrefectureModel()
				mockUseCase.EXPECT().GetPrefectureByCode(gomock.Any(), "01").Return(prefecture, nil)
			},
			wantStatus: http.StatusOK,
			wantBody: func() string {
				response := expectedGetPrefectureResponse()
				responseJSON, _ := json.Marshal(response)
				return string(responseJSON)
			},
		},
		{
			name: "Not Found",
			code: "999",
			mockSetup: func(mockUseCase *mockusecase.MockPrefectureUseCase) {
				mockUseCase.EXPECT().GetPrefectureByCode(gomock.Any(), "999").Return(nil, errors.New("not found"))
			},
			wantStatus: http.StatusNotFound,
			wantBody:   nil,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			appLogger := logger.New(logger.DefaultConfig())

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			req := httptest.NewRequest(http.MethodGet, "/prefectures/"+tt.code, nil)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			c.Params = []gin.Param{
				{
					Key:   "code",
					Value: tt.code,
				},
			}

			uc := mockusecase.NewMockPrefectureUseCase(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(uc)
			}

			mockHandler := handler.NewPrefectureHandler(appLogger, uc)
			mockHandler.GetPrefecture(c)

			a.Equal(tt.wantStatus, rec.Code)

			if tt.wantBody != nil {
				wantBody := tt.wantBody()
				resBody := rec.Body.String()
				if !cmp.Equal(wantBody, resBody) {
					t.Errorf("diff: %s", cmp.Diff(wantBody, resBody))
				}
			}
		})
	}
}

func expectedPrefectureModel() *model.Prefecture {
	return &model.Prefecture{
		ID:   1,
		Name: "北海道",
		Code: "01",
		Municipalities: []model.Municipality{
			{
				ID:                    1,
				PrefectureCode:        "01",
				OrganizationCode:      "011002",
				PrefectureNameKanji:   "北海道",
				MunicipalityNameKanji: "札幌市",
				PrefectureNameKana:    "ﾎｯｶｲﾄﾞｳ",
				MunicipalityNameKana:  "ｻｯﾎﾟﾛｼ",
				IsActive:              true,
			},
			{
				ID:                    2,
				PrefectureCode:        "13",
				OrganizationCode:      "131016",
				PrefectureNameKanji:   "東京都",
				MunicipalityNameKanji: "千代田区",
				PrefectureNameKana:    "ﾄｳｷｮｳﾄ",
				MunicipalityNameKana:  "ﾁﾖﾀﾞｸ",
				IsActive:              true,
			},
		},
	}
}

func expectedGetPrefectureResponse() *handler.GetPrefectureResponse {
	return &handler.GetPrefectureResponse{
		ID:   1,
		Name: "北海道",
		Municipalities: []*handler.Municipality{
			{
				ID:                    1,
				PrefectureCode:        "01",
				OrganizationCode:      "011002",
				PrefectureNameKanji:   "北海道",
				MunicipalityNameKanji: "札幌市",
				PrefectureNameKana:    "ﾎｯｶｲﾄﾞｳ",
				MunicipalityNameKana:  "ｻｯﾎﾟﾛｼ",
				IsActive:              true,
			},
			{
				ID:                    2,
				PrefectureCode:        "13",
				OrganizationCode:      "131016",
				PrefectureNameKanji:   "東京都",
				MunicipalityNameKanji: "千代田区",
				PrefectureNameKana:    "ﾄｳｷｮｳﾄ",
				MunicipalityNameKana:  "ﾁﾖﾀﾞｸ",
				IsActive:              true,
			},
		},
	}
}
