package brand

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BintangAldian17/voucher-redemption-service/types"
	"github.com/gorilla/mux"
)

func TestBrandServiceHandlers(t *testing.T) {
	brandStore := &mockBrandStore{}
	handler := NewHandler(brandStore)

	t.Run("it should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.CreateBrandPayload{
			Name:        "",
			Description: "",
		}

		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/brand", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/brand", handler.handleCreateBrand)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("it should success to create new brand", func(t *testing.T) {
		payload := types.CreateBrandPayload{
			Name: "brand name",
		}

		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/brand", bytes.NewBuffer(marshalled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/brand", handler.handleCreateBrand)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

}

type mockBrandStore struct{}

func (m *mockBrandStore) GetBrandByName(name string) (*types.Brand, error) {
	return nil, fmt.Errorf("brand not found")
}
func (m *mockBrandStore) CreateBrand(brand types.Brand) error {
	return nil
}
