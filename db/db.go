package db

import (
	"database/sql"
	"github.com/Gretamass/kys-backend/user"
	"log"
	_ "modernc.org/sqlite"
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

func (d *DB) AddUser(user user.User) error {
	row, err := d.db.Prepare("INSERT INTO users(email, password) VALUES (?, ?)")

	if err != nil {
		return err
	}

	_, err = row.Exec(user.Email, user.Password)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (d *DB) DeleteUser(userId int) error {
	row, err := d.db.Prepare("DELETE FROM users WHERE id=?")

	if err != nil {
		return err
	}

	_, err = row.Exec(userId)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
