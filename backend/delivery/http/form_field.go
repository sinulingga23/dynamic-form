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

type formFieldHttp struct {
	pFormFieldUsecase usecase.IPFormFieldUsecase
}

func NewFormFieldHttp(
	pFormFieldUsecase usecase.IPFormFieldUsecase,
) formFieldHttp {
	return formFieldHttp{pFormFieldUsecase: pFormFieldUsecase}
}

func (f *formFieldHttp) ServeHandler(r *chi.Mux) {
	r.Post("/api/v1/form-fields", f.HandleAddFormField)
}

func (f *formFieldHttp) HandleAddFormField(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bytesBody, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	requestData := struct {
		Data payload.PFormFieldRequest `json:"data"`
	}{}

	if errUnmarshal := json.Unmarshal(bytesBody, &requestData); errUnmarshal != nil {
		log.Println("errUnmarshal:", errUnmarshal)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := f.pFormFieldUsecase.AddFormField(r.Context(), requestData.Data)
	bytesReponse, errMarshal := json.Marshal(response)
	if errMarshal != nil {
		log.Printf("errMarshal:%v", errMarshal)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(response.StatusCode)
	w.Write(bytesReponse)
	return
}
