package server

import (
	"Yakudza/pkg/database/models"
	"Yakudza/pkg/logger"
	"Yakudza/pkg/utilities"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
)

// GetLinks - список всех ссылок
func (route Router) GetLinks(w http.ResponseWriter, r *http.Request) {
	links, err := models.AllLinks()
	if err != nil {
		logger.Error("Ошибка при получении линков: %v", err)
		SetHTTPError(w, "Ошибка на стороне сервера при попытке получения списка ссылок", http.StatusInternalServerError)
		return
	}

	sort.Slice(links, func(i, j int) bool {
		return links[i].Position < links[j].Position
	})

	str := utilities.ToJSON(links)

	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// GetLink - получение ссылки по ID
func (route Router) GetLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)

	if id <= 0 {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	link := &models.Links{ID: id}

	if err := link.FindID(); err != nil {
		SetHTTPError(w, "Ошибка при поиске ссылки", http.StatusInternalServerError)
		return
	}

	if link.Link == "" {
		SetHTTPError(w, "Ссылка не найдена", http.StatusNotFound)
		return
	}

	str := utilities.ToJSON(link)

	_, err := w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// CreateLink - создание новой ссылки
func (route Router) CreateLink(w http.ResponseWriter, r *http.Request) {
	newLink := new(models.Links)

	if err := json.NewDecoder(r.Body).Decode(newLink); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	if err := newLink.Create(); err != nil {
		SetHTTPError(w, "Ошибка при создании записи", http.StatusBadRequest)
		return
	}

	str := utilities.ToJSON(newLink)
	_, err := w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}

// DeleteLinkByID - удаление ссылки по ID
func (route Router) DeleteLinkByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["id"]
	id := utilities.StrToUint(vars)
	if id <= 0 {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	link := &models.Links{ID: id}

	if err := link.DeleteByID(); err != nil {
		SetHTTPError(w, "Ошибка при удалении ссылки", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateLink - обновление ссылки
func (route Router) UpdateLink(w http.ResponseWriter, r *http.Request) {
	id := utilities.StrToUint(mux.Vars(r)["id"])

	if id <= 0 {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	link := &models.Links{ID: id}

	if err := link.FindID(); err != nil {
		SetHTTPError(w, "Ошибка при поиске ссылки", http.StatusInternalServerError)
		return
	}

	if link.Link == "" {
		SetHTTPError(w, "Ссылка не найдена", http.StatusNotFound)
		return
	}

	logger.Info("Найденная ссылка: %s", utilities.ToJSON(link))

	updateLink := new(models.Links)

	if err := json.NewDecoder(r.Body).Decode(updateLink); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	updateLink.ID = id

	logger.Info("Новая ссылка: %s", utilities.ToJSON(updateLink))

	if err := updateLink.Update(); err != nil {
		SetHTTPError(w, "Ошибки при обновлении ссылки", http.StatusInternalServerError)
		return
	}

	link = &models.Links{ID: id}

	if err := link.FindID(); err != nil {
		SetHTTPError(w, "Ошибка при поиске ссылки", http.StatusInternalServerError)
		return
	}

	str := utilities.ToJSON(link)

	_, err := w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}
