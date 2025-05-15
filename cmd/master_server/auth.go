package masterserver

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"tic-tac-toe/cmd/helpers"
	"time"
    "net"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
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
    var userID int
    row := db.QueryRow("SELECT id FROM users WHERE email = $1", user.Email)
    err = row.Scan(&userID)
    if err != nil {
        helpers.SendErrorResponse(w, "Error scanning rows", err, http.StatusInternalServerError)
        return
    }

    ipAddr := r.RemoteAddr
    host, _, err := net.SplitHostPort(ipAddr)
    if err != nil {
        host = ""
    }
    user_agent := r.UserAgent()
    sessionID := uuid.New().String()
    _, err = db.Exec("INSERT INTO sessions (session_id, user_id, created_at, expires_at, ip_address, user_agent) VALUES ($1, $2, $3, $4, $5, $6)",sessionID, userID, time.Now(), time.Now().Add(1*time.Hour), host, user_agent)
    if err != nil {
        helpers.SendErrorResponse(w, "Error creating a session", err, http.StatusInternalServerError)
        return
    }
    http.SetCookie(w, &http.Cookie{
        Name: "session_id",
        Value:    sessionID,
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    })
    helpers.SendSuccessResponse(w, "Successfully created a new user", user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        helpers.SendErrorResponse(w, "Invalid http method", nil, http.StatusMethodNotAllowed)
    } else {
        var input struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }

        err := json.NewDecoder(r.Body).Decode(&input)
        if err != nil {
            helpers.SendErrorResponse(w, "Error decoding request", err, http.StatusBadRequest)
            return
        }

        if input.Username == "" || input.Password == "" {
            helpers.SendErrorResponse(w, "All fields are required", nil, http.StatusBadRequest)
            return
        }

        var hashedPass string
        var userID int
        row := db.QueryRow("SELECT id, password_hash FROM users WHERE username = $1", input.Username)
        err = row.Scan(&userID, &hashedPass)
        if err != nil {
            helpers.SendErrorResponse(w, "Error scanning rows", err, http.StatusInternalServerError)
            return
        }
        err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(input.Password))
        if err != nil {
            helpers.SendErrorResponse(w, "Wrong password", err, http.StatusBadRequest)
        } else {
            ipAddr := r.RemoteAddr
            host, _, err := net.SplitHostPort(ipAddr)
            if err != nil {
                ipAddr = ""
            }
            user_agent := r.UserAgent()
            sessionID := uuid.New().String()
            _, err = db.Exec("INSERT INTO sessions (session_id, user_id, created_at, expires_at, ip_address, user_agent) VALUES ($1, $2, $3, $4, $5, $6)",sessionID, userID, time.Now(), time.Now().Add(1*time.Hour), host, user_agent)
            if err != nil {
                helpers.SendErrorResponse(w, "Error creating a session", err, http.StatusInternalServerError)
                return
            }
            http.SetCookie(w, &http.Cookie{
                Name: "session_id",
                Value:    sessionID,
                Expires:  time.Now().Add(24 * time.Hour),
                HttpOnly: true,
                Secure:   true,
                SameSite: http.SameSiteLaxMode,
            })
            helpers.SendSuccessResponse(w, "Login successfull", nil)
        }
    }
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        helpers.SendErrorResponse(w, "Invalid http method", nil, http.StatusMethodNotAllowed)
    } else {
        cookie, err := r.Cookie("session_id")
        if err == nil {
        _, _ = db.Exec("DELETE FROM sessions WHERE session_id = $1", cookie.Value)
        }
        http.SetCookie(w, &http.Cookie{
            Name:     "session_id",
            Value:    "",
            Expires:  time.Unix(0, 0),
            MaxAge:   -1,
            HttpOnly: true,
            Secure:   true,
            SameSite: http.SameSiteLaxMode,
        })
        helpers.SendSuccessResponse(w, "User successfully logged out", nil)
    }
}
