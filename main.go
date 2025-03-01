package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

var db *sqlx.DB

type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func GetUsers() ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT id, name, age FROM users")
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(id int) (*User, error) {
	query := "SELECT id, name, age FROM users WHERE id=:id"
	params := map[string]interface{}{"id": id}

	var user User
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Get(&user, params)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func InsertUser(name string, age int) (int64, error) {
	query := "INSERT INTO users (name, age) OUTPUT INSERTED.id VALUES (:name, :age)"
	params := map[string]interface{}{
		"name": name,
		"age":  age,
	}

	var id int64
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.Get(&id, params)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateUser(id int, name string, age int) error {
	query := "UPDATE users SET name=:name, age=:age WHERE id=:id"
	params := map[string]interface{}{
		"id":   id,
		"name": name,
		"age":  age,
	}

	_, err := db.NamedExec(query, params)
	return err
}

func DeleteUser(id int) error {
	query := "DELETE FROM dbo.users WHERE id = :id"
	_, err := db.NamedExec(query, map[string]interface{}{"id": id})
	return err
}

func main() {
	// Configure slog with a text handler
	slog.SetDefault(slog.New(slog.NewTextHandler(log.Writer(), nil)))

	// Database connection
	dsn := "sqlserver://sa:Test1234@localhost:1433?database=fiet&encrypt=disable"
	var err error
	db, err = sqlx.Connect("sqlserver", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	fmt.Println("Connected to MSSQL successfully!")

	// Fetch users
	users, err := GetUsers()
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}

	// Log users with slog
	slog.Info("Fetched users", "users", users)

	// Insert a new user
	// userid, err := InsertUser("Alice", 40)
	// if err != nil {
	// 	log.Fatalf("Failed to insert user: %v", err)
	// }

	// Log the new user ID
	// slog.Info("Inserted user", "id", userid)

	// Fetch a user by ID
	user, err := GetUserByID(2)
	if err != nil {
		log.Fatalf("Failed to get user by ID: %v", err)
	}

	// Log the user with slog
	slog.Info("Fetched user by ID", "user", user)

	// Update a user
	// err = UpdateUser(1, "Bass", 30)
	// if err != nil {
	// 	log.Fatalf("Failed to update user: %v", err)
	// }

	// Log the updated user
	// slog.Info("Updated user", "id", 1)

	// Delete a user
	// err = DeleteUser(1)
	// if err != nil {
	// 	log.Fatalf("Failed to delete user: %v", err)
	// }

	// // Log the deleted user
	// slog.Info("Deleted user", "id", 1)

	// Correct query to retrieve a value
	var names []string
	err = db.Select(&names, "SELECT name FROM sys.databases")
	if err != nil {
		log.Fatalf("Failed to get names: %v", err)
	}

	// Print the database names
	fmt.Println("Databases:", names)
}
