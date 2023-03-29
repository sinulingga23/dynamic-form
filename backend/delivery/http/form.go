package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sinulingga23/dynamic-form/backend/api/usecase"
	"github.com/sinulingga23/dynamic-form/backend/payload"
)

type formHttp struct {
	pFormUsecase usecase.IPFormUsecase
}

func NewFormHttp(
	pFormUsecase usecase.IPFormUsecase,
) formHttp {
	return formHttp{pFormUsecase: pFormUsecase}
}

func (f *formHttp) ServeHandler(r *chi.Mux) {
	r.Post("/api/v1/forms", f.HandleAddPForm)
	r.Get("/api/v1/forms/{partnerId}/partner", f.HandleGetPFormsByPPartnerId)
	r.Get("/api/v1/forms/{id}", f.HandleGetPFormById)
}

func (f *formHttp) HandleAddPForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bytesBody, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestData := struct {
		Data payload.PFormRequest `json:"data"`
	}{}

	if errUnmarshal := json.Unmarshal(bytesBody, &requestData); errUnmarshal != nil {
		log.Println("errUnmarshal:", errUnmarshal)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := f.pFormUsecase.AddPForm(r.Context(), requestData.Data)
	bytesReponse, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		log.Printf("errMarshal:%v", errMarshal)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(response.StatusCode)
	w.Write(bytesReponse)
	return
}

func (f *formHttp) HandleGetPFormsByPPartnerId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	partnerId := chi.URLParam(r, "partnerId")
	response := f.pFormUsecase.GetPFormsByPPartnerId(r.Context(), partnerId)

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

func (f *formHttp) HandleGetPFormById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	response := f.pFormUsecase.GetPFormById(r.Context(), id)

	bytesResponse, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		log.Printf("errMashal: %v", errMarshal)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(response.StatusCode)
	w.Write(bytesResponse)
	return
}
