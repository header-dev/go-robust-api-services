package main

import (
	"header-dev/go-robust-api-services/todo"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDatabase(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	return nil
}

func createTodosTable() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS todos (
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
        "title" TEXT
    );`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	dbPath := "test.db"
	if len(dbPath) == 0 {
		panic("specify the SQLITE_DB_PATH environment variable")
	}
	err := initDatabase(dbPath)

	if err != nil {
		panic("failed to connect database")
	}

	err = db.Ping()
	if err != nil {
		panic("error initializing DB connection: ping error: ")
	}

	//Auto migration
	createTodosTable()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	handler := todo.NewTodoHandler(db)
	r.POST("/todos", handler.NewTask)

	r.Run()
}
