package db

import (
	"database/sql"
	"fmt"
	"github.com/Gretamass/kys-backend/user"
	_ "modernc.org/sqlite"
	"strings"
)

type DB struct {
	db *sql.DB
}

func ConnectDatabase() (*DB, error) {
	db, err := sql.Open("sqlite", "./sqlite.db")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

// USER methods

func (d *DB) GetUsers() ([]user.User, error) {
	rows, err := d.db.Query("SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]user.User, 0)

	for rows.Next() {
		singleUser := user.User{}
		err = rows.Scan(&singleUser.Id, &singleUser.Email, &singleUser.Password, &singleUser.CreatedAt)

		if err != nil {
			return nil, err
		}

		users = append(users, singleUser)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (d *DB) GetUserById(userId int) (user.User, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE id = ?", userId)

	singleUser := user.User{}
	err := row.Scan(&singleUser.Id, &singleUser.Email, &singleUser.Password, &singleUser.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, fmt.Errorf("no rows found with id %d", userId)
		}
		return user.User{}, err
	}

	return singleUser, nil
}

func (d *DB) AddUser(user user.User) error {
	row, err := d.db.Prepare("INSERT INTO users (email, password) VALUES (?, ?)")

	if err != nil {
		return err
	}

	_, err = row.Exec(user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) UpdateUser(userId int, request user.User) error {
	query := "UPDATE users SET "
	var args []interface{}

	if request.Email != "" {
		query += "email = ?, "
		args = append(args, request.Email)
	}

	if request.Password != "" {
		query += "password = ?, "
		args = append(args, request.Password)
	}

	query = strings.TrimRight(query, ", ")
	query += " WHERE id = ?"
	args = append(args, userId)

	row, err := d.db.Prepare(query)

	if err != nil {
		return err
	}

	_, err = row.Exec(args...)

	if err != nil {
		return err
	}

	return nil
}

func (d *DB) DeleteUser(userId int) error {
	row, err := d.db.Prepare("DELETE FROM users WHERE id = ?")

	if err != nil {
		return err
	}

	result, err := row.Exec(userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows found with id %d", userId)
	}

	return nil
}

// ADMIN methods

func (d *DB) GetAdmins() ([]user.Admin, error) {
	rows, err := d.db.Query("SELECT * FROM admins")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	admins := make([]user.Admin, 0)

	for rows.Next() {
		singleAdmin := user.Admin{}
		err = rows.Scan(&singleAdmin.Id, &singleAdmin.Email, &singleAdmin.Password)

		if err != nil {
			return nil, err
		}

		admins = append(admins, singleAdmin)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (d *DB) AddAdmin(admin user.Admin) error {
	row, err := d.db.Prepare("INSERT INTO admins (email, password) VALUES (?, ?)")

	if err != nil {
		return err
	}

	_, err = row.Exec(admin.Email, admin.Password)
	if err != nil {
		return err
	}

	return nil
}
