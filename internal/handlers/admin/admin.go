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

// GET /admin/users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := models.GetAllUsers()
    if err != nil {
        http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// POST /admin/users
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

// PUT /admin/users/{id}
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

// DELETE /admin/users/{id}
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