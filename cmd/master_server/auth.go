package masterserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tic-tac-toe/cmd/helpers"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func HashPassword() string {

}
func Register(w http.ResponseWriter, r *http.Request, db *helpers.Database) {
	if r.Method != "POST" {
		helpers.SendErrorResponse(w, "Invalid method", nil, http.StatusMethodNotAllowed)
	} else {
		var input struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			helpers.SendErrorResponse(w, "Error decoding request", err, http.StatusInternalServerError)
		}

		hashedPass := HashPassword()
		current_time := time.Now()

		var user = User{
			ID:           uuid.New(),
			Username:     input.Username,
			Email:        input.Email,
			PasswordHash: hashedPass,
			CreatedAt:    current_time,
			UpdatedAt:    current_time,
		}

		_, err = db.Conn.Exec("INSERT INTO users (id, username, email, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)", user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			helpers.SendErrorResponse(w, "Error registering user", err, http.StatusInternalServerError)
		}

		helpers.SendSuccessResponse(w, "Successfully created a new user", user)

	}
}
func Login() {
	fmt.Println("Hello auth")
}
func Logout() {
	fmt.Println("Hello auth")
}
