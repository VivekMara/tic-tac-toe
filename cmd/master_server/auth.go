package masterserver

import (
	"encoding/json"
	"database/sql"
	"fmt"
	"net/http"
	"tic-tac-toe/cmd/helpers"
	"time"
	"strings"

	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

type User struct {
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var db *sql.DB

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        helpers.SendErrorResponse(w, "Invalid method", nil, http.StatusMethodNotAllowed)
        return
    }

    var input struct {
        Username string `json:"username"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        helpers.SendErrorResponse(w, "Error decoding request", err, http.StatusBadRequest)
        return
    }

    if input.Username == "" || input.Email == "" || input.Password == "" {
        helpers.SendErrorResponse(w, "All fields are required", nil, http.StatusBadRequest)
        return
    }

    hashedPass, err := hash(input.Password)
    if err != nil {
        helpers.SendErrorResponse(w, "Error hashing the password", err, http.StatusInternalServerError)
        return
    }

    now := time.Now()
    user := User{
        Username:     input.Username,
        Email:        input.Email,
        PasswordHash: hashedPass,
        CreatedAt:    now,
        UpdatedAt:    now,
    }

    _, err = db.Exec(`INSERT INTO users 
        (username, email, password_hash, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5)`,
        user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
    )
    if err != nil {
        if strings.Contains(err.Error(), "duplicate key") {
            helpers.SendErrorResponse(w, "User already exists", err, http.StatusConflict)
        } else {
            helpers.SendErrorResponse(w, "Error registering user", err, http.StatusInternalServerError)
        }
        return
    }

    helpers.SendSuccessResponse(w, "Successfully created a new user", user)
}

func Login() {
	fmt.Println("Hello auth")
}

func Logout() {
	fmt.Println("Hello auth")
}
