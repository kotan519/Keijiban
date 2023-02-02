package dbrepo

import (
	"context"
	"errors"
	//"fmt"
	"log"
	"time"

	"github.com/kotan519/keijiban/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m postgresDBRepo) InsertData(as models.TokumeiPostData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into threads_table (name, body, created_at, updated_at)
			values ($1, $2, $3, $4)`

	_, err := m.DB.ExecContext(ctx, stmt,
		as.Title,
		as.Text,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetData() (*models.TokumeiPostData, error) {
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select id name body created_at FROM threads_table`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := &models.TokumeiPostData{}
	for rows.Next() {

		err := rows.Scan(&res.ID, &res.Title, &res.Text, &res.CreatedAt)
		if err != nil {
			return nil, err
		}
		//fmt.Println(fmt.Sprintf("ID: %d, Title: %s, Text: %s, Created_At: %s"), res.ID, res.Title, res.Text, res.CreatedAt)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return res, nil
}

func (m *postgresDBRepo) GetThreadList() ([]models.TokumeiPostData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var threadsdata []models.TokumeiPostData

	query := `
		SELECT r.id, r.name, r.body, r.created_at, r.updated_at
		FROM threads_table r
		ORDER BY r.id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return threadsdata, err
	}

	defer rows.Close()

	for rows.Next() {
		var i models.TokumeiPostData
		err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Text,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return threadsdata, err
		}
		threadsdata = append(threadsdata, i)
	}

	if err := rows.Err(); err != nil {
		return threadsdata, err
	}

	return threadsdata, nil
}

func (m *postgresDBRepo) GetThreadData(num int) ([]models.TokumeiPostData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var threadsdata []models.TokumeiPostData

	query := `
		SELECT r.id, r.name, r.body, r.created_at, r.updated_at
		FROM threads_table r
		WHERE r.id = $1
		ORDER BY r.id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, num)
	if err != nil {
		return threadsdata, nil
	}

	defer rows.Close()

	for rows.Next() {
		var i models.TokumeiPostData
		err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Text,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return threadsdata, err
		}
		threadsdata = append(threadsdata, i)
	}

	if err := rows.Err(); err != nil {
		return threadsdata, err
	}

	return threadsdata, nil
}

func (m postgresDBRepo) InsertCommentData(as models.TokumeiPostData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into comments_table (threads_id, name, body, created_at, updated_at)
			values ($1, $2, $3, $4, $5)`

	_, err := m.DB.ExecContext(ctx, stmt,
		as.ThreadID,
		as.Title,
		as.Text,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetCommentData(num int) ([]models.TokumeiPostData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var commentsdata []models.TokumeiPostData

	query := `
		SELECT r.threads_id, r.name, r.body, r.created_at, r.updated_at
		FROM comments_table r
		WHERE r.threads_id = $1
		ORDER BY r.id ASC
	`

	rows, err := m.DB.QueryContext(ctx, query, num)
	if err != nil {
		return commentsdata, nil
	}

	defer rows.Close()

	for rows.Next() {
		var i models.TokumeiPostData
		err := rows.Scan(
			&i.ThreadID,
			&i.Title,
			&i.Text,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return commentsdata, err
		}
		commentsdata = append(commentsdata, i)
	}

	if err := rows.Err(); err != nil {
		return commentsdata, err
	}

	return commentsdata, nil
}

func (m *postgresDBRepo) Authenticate(email, testpassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testpassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("パスワードが違います")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, err
}

func (m postgresDBRepo) InsertUserData(as models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into users (username, email, password, created_at, updated_at)
			values ($1, $2, $3, $4, $5)`

	_, err := m.DB.ExecContext(ctx, stmt,
		as.UserName,
		as.Email,
		as.Password,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

