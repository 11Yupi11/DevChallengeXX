package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dev-challenge/internal/models"
	mock_services "dev-challenge/internal/services/mock"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	mockLogger = &logrus.Logger{}
)

func TestHandler_getValue(t *testing.T) {
	type mockBehavior func(r *mock_services.MockExcelLikeService)

	type Test struct {
		Name                 string
		url                  string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}
	tests := [...]Test{
		{
			Name: "Valid sheet and cell ID",
			url:  "/api/v1/sheetID1/cellID1",
			mockBehavior: func(r *mock_services.MockExcelLikeService) {
				r.EXPECT().GetCellInput(gomock.Any(), strings.ToLower("sheetID1"), strings.ToLower("cellID1")).Return(&models.Data{
					Value: "SomeValue",
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"value\":\"SomeValue\",\"result\":\"\"}\n",
		},
		{
			Name: "Value not found",
			url:  "/api/v1/sheetID2/cellID2",
			mockBehavior: func(r *mock_services.MockExcelLikeService) {
				r.EXPECT().GetCellInput(gomock.Any(), strings.ToLower("sheetID2"), strings.ToLower("cellID2")).Return(&models.Data{
					Value: "",
				}, nil)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "{\"code\":\"404\",\"message\":\"value not found\"}\n",
		},
		{
			Name:                 "Invalid sheet ID",
			url:                  "/api/v1/|heetID/cell1",
			mockBehavior:         func(r *mock_services.MockExcelLikeService) {}, // No expected calls on the mock in this case
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "{\"code\":\"404\",\"message\":\"not correct params\"}\n",
		},
		{
			Name:                 "Invalid cell ID",
			url:                  "/api/v1/sheetID1/^cell1",
			mockBehavior:         func(r *mock_services.MockExcelLikeService) {}, // No expected calls on the mock in this case
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "{\"code\":\"404\",\"message\":\"not correct params\"}\n",
		},
		{
			Name: "Store not responded",
			url:  "/api/v1/sheetID1/cellID1",
			mockBehavior: func(r *mock_services.MockExcelLikeService) {
				r.EXPECT().GetCellInput(gomock.Any(), strings.ToLower("sheetID1"), strings.ToLower("cellID1")).Return(nil, errors.New("store not responded"))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "{\"code\":\"404\",\"message\":\"store not responded\"}\n",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_services.NewMockExcelLikeService(ctrl)
			test.mockBehavior(m)

			r := chi.NewRouter()
			h := &ExcelLikeHandler{
				ELS: m,
				Log: mockLogger,
			}
			r.HandleFunc("/api/v1/{sheet_id}/{cell_id}", h.getValue)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", test.url, bytes.NewBuffer(nil))
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_addValue(t *testing.T) {
	type mockBehavior func(r *mock_services.MockExcelLikeService)

	type Test struct {
		Name                 string
		url                  string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}
	tests := [...]Test{
		{
			Name:      "Successful add value",
			url:       "/api/v1/sheetID1/cellID1",
			inputBody: `{"value": "1"}`,
			mockBehavior: func(r *mock_services.MockExcelLikeService) {
				r.EXPECT().AddCellInputTX(gomock.Any(), strings.ToLower("sheetID1"), strings.ToLower("cellID1"), &models.Data{Value: "1"}).Return(&models.Data{
					Value:  "1",
					Result: "1.000000",
				}, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: "{\"value\":\"1\",\"result\":\"1.000000\"}\n",
		},
		{
			Name:      "Value updated",
			url:       "/api/v1/sheetID1/cellID1",
			inputBody: `{"value": "2"}`,
			mockBehavior: func(r *mock_services.MockExcelLikeService) {
				r.EXPECT().AddCellInputTX(gomock.Any(), strings.ToLower("sheetID1"), strings.ToLower("cellID1"), &models.Data{Value: "2"}).Return(&models.Data{
					Value:  "2",
					Result: "2.000000",
				}, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: "{\"value\":\"2\",\"result\":\"2.000000\"}\n",
		},
		{
			Name:                 "Invalid sheet ID",
			url:                  "/api/v1/|heetID/cell1",
			mockBehavior:         func(r *mock_services.MockExcelLikeService) {}, // No expected calls on the mock in this case
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedResponseBody: "{\"code\":\"422\",\"message\":\"not correct params\"}\n",
		},
		{
			Name:                 "Invalid cell ID",
			url:                  "/api/v1/sheetID1/^cell1",
			mockBehavior:         func(r *mock_services.MockExcelLikeService) {}, // No expected calls on the mock in this case
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedResponseBody: "{\"code\":\"422\",\"message\":\"not correct params\"}\n",
		},
		{
			Name:      "Mistake in adding",
			url:       "/api/v1/sheetID1/cellID1",
			inputBody: `{"value": "2"}`,
			mockBehavior: func(r *mock_services.MockExcelLikeService) {
				r.EXPECT().AddCellInputTX(gomock.Any(), strings.ToLower("sheetID1"), strings.ToLower("cellID1"), &models.Data{Value: "2"}).Return(nil, errors.New("store not responded"))
			},
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedResponseBody: "{\"value\":\"2\",\"result\":\"ERROR\"}\n",
		},
		{
			Name:                 "Not correct input",
			url:                  "/api/v1/sheetID1/cellID1",
			inputBody:            `{123}`,
			mockBehavior:         func(r *mock_services.MockExcelLikeService) {},
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedResponseBody: "{\"code\":\"422\",\"message\":\"can't unmarshal request body\"}\n",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_services.NewMockExcelLikeService(ctrl)
			test.mockBehavior(m)

			r := chi.NewRouter()
			h := &ExcelLikeHandler{
				ELS: m,
				Log: mockLogger,
			}
			r.HandleFunc("/api/v1/{sheet_id}/{cell_id}", h.addValue)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", test.url, strings.NewReader(test.inputBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getAllValues(t *testing.T) {
	type mockBehavior func(r *mock_services.MockExcelLikeService)

	type Test struct {
		Name                 string
		url                  string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}
	tests := [...]Test{
		{
			Name: "Valid sheet ID",
			url:  "/api/v1/sheetID1",
			mockBehavior: func(r *mock_services.MockExcelLikeService) {
				data := map[string]models.Data{
					"cell1": {
						Value:  "1",
						Result: "1.000000",
					},
				}
				r.EXPECT().GetSheetInput(gomock.Any(), strings.ToLower("sheetID1")).Return(data, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "{\"cell1\":{\"value\":\"1\",\"result\":\"1.000000\"}}\n",
		},
		{
			Name: "Value not found",
			url:  "/api/v1/sheetID2",
			mockBehavior: func(r *mock_services.MockExcelLikeService) {
				r.EXPECT().GetSheetInput(gomock.Any(), strings.ToLower("sheetID2")).Return(nil, nil)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "{\"code\":\"404\",\"message\":\"value not found\"}\n",
		},
		{
			Name:                 "Invalid sheet ID",
			url:                  "/api/v1/|heetID",
			mockBehavior:         func(r *mock_services.MockExcelLikeService) {}, // No expected calls on the mock in this case
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "{\"code\":\"404\",\"message\":\"not correct params\"}\n",
		},
		{
			Name: "Store not responded",
			url:  "/api/v1/sheetID1",
			mockBehavior: func(r *mock_services.MockExcelLikeService) {
				r.EXPECT().GetSheetInput(gomock.Any(), strings.ToLower("sheetID1")).Return(map[string]models.Data{}, errors.New("store not responded"))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: "{\"code\":\"404\",\"message\":\"store not responded\"}\n",
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_services.NewMockExcelLikeService(ctrl)
			test.mockBehavior(m)

			r := chi.NewRouter()
			h := &ExcelLikeHandler{
				ELS: m,
				Log: mockLogger,
			}
			r.HandleFunc("/api/v1/{sheet_id}", h.getAllValues)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", test.url, bytes.NewBuffer(nil))
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
