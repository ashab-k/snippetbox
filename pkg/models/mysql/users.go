package mysql

import (
	"database/sql"
	"strings"

	"github.com/ashab-k/snippetbox/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
    DB *sql.DB
}

// We'll use the Insert method to add a new record to the users table.
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }

    stmt := `INSERT INTO users (name, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())`

    // Use the Exec() method to insert the user details and hashed password
    // into the users table. If this returns an error, we try to type assert
    // it to a *mysql.MySQLError object so we can check if the error number is
    // 1062 and, if it is, we also check whether or not the error relates to
    // our users_uc_email key by checking the contents of the message string.
    // If it does, we return an ErrDuplicateEmail error. Otherwise, we just
    // return the original error (or nil if everything worked).
    _, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
    if err != nil {
        if mysqlErr, ok := err.(*mysql.MySQLError); ok {
            if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message, "users_uc_email") {
                return models.ErrDuplicateEmail
            }
        }
    }
    return err
}

// We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := m.DB.QueryRow("SELECT id , hashed_password from users WHERE email = ?" ,email)

	err := row.Scan(&id , &hashedPassword)
	if err == sql.ErrNoRows{
		return 0 , models.ErrInvalidCredentials
	}else if err != nil {
		return 0 , err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword , []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword{
		return 0 , models.ErrInvalidCredentials
	} else if err != nil {
		return 0 , err
	}


    return id, nil
}

// We'll use the Get method to fetch details for a specific user based
// on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
    s := &models.User{}

    stmt := `SELECT id, name, email, created FROM users WHERE id = ?`
    err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created)
    if err == sql.ErrNoRows {
        return nil, models.ErrNoRecord
    } else if err != nil {
        return nil, err
    }

    return s, nil
}