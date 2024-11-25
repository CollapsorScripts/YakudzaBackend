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

// GetLinks godoc
// @Summary      Список всех ссылок
// @Description  Массив с ссылками в базе данных
// @Tags         Links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  models.Links
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /links [get]
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

// GetLink godoc
// @Summary      Поиск ссылки
// @Description  Поиск ссылки по ID
// @Tags         Links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Link ID"
// @Success      200  {object}  models.Links
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /links/{id} [get]
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

// CreateLink godoc
// @Summary      Создание ссылки
// @Description  Создание новой сущности ссылки
// @Tags         Links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        link body models.Links false "Сущность ссылки"
// @Success      200  {object}  models.Links
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /links [post]
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

// DeleteLinkByID godoc
// @Summary      Удаление ссылки
// @Description  Удаление ссылки по ID
// @Tags         Links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "ID ссылки"
// @Success      200              {string}  string    "ok"
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /links/{id} [delete]
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

// UpdateLink godoc
// @Summary      Обновление ссылки
// @Description  Обновление сущности ссылки
// @Tags         Links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID ссылки"
// @Param        link body models.Links true "Модель для обновления"
// @Success      200  {object}  models.Links
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /links/{id} [put]
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
