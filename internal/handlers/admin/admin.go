package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "go-auth-app/internal/models"
    "go-auth-app/internal/utils"
    "go-auth-app/internal/config"

)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Добро пожаловать, админ!")
}

// GetAllUsers godoc
// @Summary Получить всех пользователей
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.User
// @Router /admin/users [get]
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := models.GetAllUsers()
    if err != nil {
        http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// CreateUser godoc
// @Summary      Создать пользователя
// @Description  Доступно только администратору
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        input body models.User true "Новый пользователь"
// @Security     BearerAuth
// @Success      201 {string} string "Пользователь создан"
// @Failure      400 {string} string "Невалидные данные"
// @Router       /admin/users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
        return
    }

    if user.Username == "" || user.Email == "" || user.Password == "" || user.Role == "" {
        http.Error(w, "Все поля обязательны", http.StatusBadRequest)
        return
    }

    if err := models.CreateUser(user); err != nil {
        http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    fmt.Fprintln(w, "Пользователь успешно создан")
}

// UpdateUser godoc
// @Summary      Обновить пользователя
// @Description  Обновляет информацию пользователя по ID
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Param        input body models.User true "Обновлённые данные"
// @Security     BearerAuth
// @Success      200 {string} string "Пользователь обновлён"
// @Failure      400 {string} string "Ошибка запроса"
// @Failure      500 {string} string "Ошибка сервера"
// @Router       /admin/users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Некорректный ID", http.StatusBadRequest)
        return
    }

    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Неверный формат данных", http.StatusBadRequest)
        return
    }

    if user.Password != "" {
        hashedPassword, err := utils.HashPassword(user.Password)
        if err != nil {
            http.Error(w, "Ошибка хеширования пароля", http.StatusInternalServerError)
            return
        }

        query := `UPDATE users
                  SET username = $1, email = $2, role = $3, password = $4, updated_at = NOW()
                  WHERE id = $5`

        _, err = config.DB.Exec(
            r.Context(),
            query,
            user.Username,
            user.Email,
            user.Role,
            hashedPassword,
            id,
        )

        if err != nil {
            fmt.Println("Ошибка при обновлении с паролем:", err)
            http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
            return
        }
    } else {
        if err := models.UpdateUserByID(id, user); err != nil {
            fmt.Println("Ошибка при обновлении без пароля:", err)
            http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
            return
        }
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Пользователь успешно обновлён")
}

// DeleteUser godoc
// @Summary      Удалить пользователя
// @Description  Удаляет пользователя по ID
// @Tags         admin
// @Param        id path int true "ID пользователя"
// @Security     BearerAuth
// @Success      200 {string} string "Пользователь удалён"
// @Failure      500 {string} string "Ошибка удаления"
// @Router       /admin/users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Некорректный ID", http.StatusBadRequest)
        return
    }

    if err := models.DeleteUserByID(id); err != nil {
        http.Error(w, "Ошибка удаления пользователя", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Пользователь успешно удалён")
}

// AdminDashboard godoc
// @Summary      Панель администратора
// @Description  Показывает статистику: количество пользователей, активные сессии, последние входы
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      500 {string} string "Ошибка сервера"
// @Router       /admin/dashboard [get]
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
    data, err := models.GetDashboardData()
    if err != nil {
        http.Error(w, "Ошибка получения данных для дашборда", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}