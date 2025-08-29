package handler

import (
	"log/slog"
	"net/http"
	"text/template"

	"1337b04rd/internal/domain"
)

func (h *handler) renderError(w http.ResponseWriter, errPage *domain.ErrorPageData) {
	if errPage == nil {
		errPage = &domain.ErrorPageData{
			Code:    http.StatusInternalServerError,
			Message: "errPage was nil",
		}
	}
	w.WriteHeader(errPage.Code)
	tmpl, err := template.ParseFiles("web/templates/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, errPage)
	if err != nil {
		slog.Error(err.Error())
	}
}
