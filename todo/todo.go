package todo

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID    int64  `json:"id"`
	Title string `json:"text"`
}

type TodoHandler struct {
	db *sql.DB
}

func NewTodoHandler(db *sql.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (t *TodoHandler) NewTask(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := "INSERT INTO todos (title) VALUES (?)"
	result, err := t.db.Exec(query, todo.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	todo.ID = id
	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.ID,
	})
}

func (t *TodoHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	var todo Todo

	query := "SELECT id, title FROM todos WHERE id = ?"
	err := t.db.QueryRow(query, id).Scan(&todo.ID, &todo.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Todo not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (t *TodoHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := "UPDATE todos SET title = ? WHERE id = ?"
	_, err := t.db.Exec(query, todo.Title, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo updated successfully",
	})
}

func (t *TodoHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM todos WHERE id = ?"
	_, err := t.db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo deleted successfully",
	})
}
