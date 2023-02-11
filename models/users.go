package models

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
}

func ConnectDatabase() error {
	db, err := sql.Open("sqlite", "./sqlite.db")
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	fmt.Println("Connected to the database.")

	DB = db
	return nil
}

func GetUsers() ([]User, error) {
	rows, err := DB.Query("SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		singleUser := User{}
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

func AddUser(c *gin.Context) error {
	var newUser User

	if err := c.BindJSON(&newUser); err != nil {
		return err
	}

	fmt.Println(c.BindJSON(&newUser))
	row, err := DB.Prepare("INSERT INTO users(email, password) VALUES (?, ?)")

	if err != nil {
		return err
	}

	_, err = row.Exec(&newUser.Email, &newUser.Password)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
