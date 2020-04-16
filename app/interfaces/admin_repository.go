package interfaces

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/bmf-san/gobel-api/app/domain"
	"github.com/bmf-san/gobel-api/app/usecases"
	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
)

// An AdminRepository is a repository for an authentication.
type AdminRepository struct {
	ConnMySQL *sql.DB
	ConnRedis *redis.Client
}

// FindIDByToken returns the entity identified by the given token.
func (ar *AdminRepository) FindIDByToken(token string) (int, error) {
	id, err := ar.ConnRedis.Get(token).Result()
	if err != nil {
		return 0, nil
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		return 0, nil
	}

	return i, nil
}

// FindByCredential returns the entity identified by the given credential.
func (ar *AdminRepository) FindByCredential(req usecases.RequestCredential) (admin domain.Admin, err error) {
	const query = `
		SELECT
			id,
			name,
			email,
			password
		FROM
			admins
		WHERE
			email = ?
	`
	row, err := ar.ConnMySQL.Query(query, req.Email)

	defer row.Close()

	if err != nil {
		return
	}

	var id int
	var name string
	var password string
	var email string
	row.Next()
	if err = row.Scan(&id, &name, &email, &password); err != nil {
		return
	}
	admin.ID = id
	admin.Name = name
	admin.Email = email
	admin.Password = password

	return
}

// SaveSessionByID saves session by the given id.
func (ar *AdminRepository) SaveSessionByID(id int) (token string, err error) {
	token = uuid.NewV4().String()
	if err := ar.ConnRedis.Set(token, strconv.Itoa(id), 3600*24*7*time.Second).Err(); err != nil {
		return "", err
	}

	return token, nil
}
