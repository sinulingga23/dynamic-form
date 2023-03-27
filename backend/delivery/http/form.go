package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sinulingga23/dynamic-form/backend/api/usecase"
)

type formHtp struct {
	pFormUsecase usecase.IPFormUsecase
}

func NewFormHttp(
	pFormUsecase usecase.IPFormUsecase,
) formHtp {
	return formHtp{pFormUsecase: pFormUsecase}
}

func (f *formHtp) HandleGetFormsByPartnerId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	partnerId := chi.URLParam(r, "partnerId")
	response := f.pFormUsecase.GetFormsByPartnerId(r.Context(), partnerId)

	bytesResponse, errMarhsal := json.Marshal(response)
	if errMarhsal != nil {
		log.Printf("errMarshal: %v", errMarhsal)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(response.StatusCode)
	w.Write(bytesResponse)
	return
}
