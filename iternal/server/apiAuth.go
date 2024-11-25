package server

import (
	"Yakudza/pkg/database/models"
	"Yakudza/pkg/logger"
	"Yakudza/pkg/token"
	"Yakudza/pkg/utilities"
	"encoding/json"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Login godoc
// @Summary      Авторизация
// @Description  Авторизация пользователя
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user body LoginRequest true "Данные для авторизации пользователя"
// @Success      200  {object}  LoginResponse
// @Failure      400  {object}  HTTPError
// @Failure      401  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /auth/login [post]
func (route Router) Login(w http.ResponseWriter, r *http.Request) {
	request := new(LoginRequest)

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		SetHTTPError(w, "Неверные аргументы", http.StatusBadRequest)
		return
	}

	if request.Login == "" || len(request.Login) == 0 {
		logger.Error("Поле login пустое")
		SetHTTPError(w, "Поле \"Login\" не может быть пустым", http.StatusBadRequest)
		return
	}

	if request.Password == "" || len(request.Password) == 0 {
		logger.Error("Поле password пустое")
		SetHTTPError(w, "Поле \"Password\" не может быть пустым", http.StatusBadRequest)
		return
	}

	accessory, err := models.ComparePassword(request.Login, request.Password)
	if err != nil {
		if strings.Contains(err.Error(), "Пользователь не найден") {
			SetHTTPError(w, err.Error(), http.StatusUnauthorized)
		} else {
			SetHTTPError(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
		}

		return
	}

	if !accessory {
		SetHTTPError(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	user := &models.User{Login: request.Login}
	if err := user.FindUserLogin(); err != nil {
		SetHTTPError(w, "Ошибка на стороне сервера", http.StatusInternalServerError)
		return
	}

	jwtToken, err := token.CreateToken(user, route.cfg)
	if err != nil {
		logger.Error("Ошибка при создании JWT токена: %v", err)
		SetHTTPError(w, "Ошибка при создании JWT токена", http.StatusInternalServerError)
		return
	}

	response := &LoginResponse{Token: jwtToken}

	str := utilities.ToJSON(response)
	_, err = w.Write([]byte(str))
	if err != nil {
		logger.Error("%s", err.Error())
	}
}
