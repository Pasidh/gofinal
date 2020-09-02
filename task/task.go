package task

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type Done struct {
	Message string `json:"message"`
}

var err error
var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

}

func CreateTable() {
	createTb := `
	CREATE TABLE IF NOT EXISTS customers(
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`

	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table ", err)
	}

	log.Println("Create table success")
}

func CreateTodosHandler(c *gin.Context) {
	m := Customer{}
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	row := db.QueryRow("INSERT INTO customers (name, email , status) values ($1, $2 ,$3)  RETURNING id", m.Name, m.Email, m.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err := row.Scan(&m.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, m)

}

func GetTodosHandler(c *gin.Context) {
	status := c.Query("status")

	stmt, err := db.Prepare("SELECT id, name ,email, status FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	customer := []Customer{}
	for rows.Next() {
		m := Customer{}

		err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		customer = append(customer, m)
	}

	mm := []Customer{}

	for _, item := range customer {
		if status != "" {
			if item.Status == status {
				mm = append(mm, item)
			}
		} else {
			mm = append(mm, item)
		}
	}

	c.JSON(http.StatusOK, mm)
}

func GetTodoByIdHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, name , email , status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	m := &Customer{}

	err = row.Scan(&m.ID, &m.Name, &m.Email, &m.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, m)
}

func UpdateTodosHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, name ,email, status FROM customers where id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)

	m := &Customer{}

	err = row.Scan(&m.ID, &m.Name, &m.Email, &m.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := c.ShouldBindJSON(m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err = db.Prepare("UPDATE customers SET name=$2, email=$3 , status=$4 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := stmt.Exec(m.ID, m.Name, m.Email, m.Status); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, m)
}

func DeleteTodosHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM customers WHERE id = $1")
	if err != nil {
		log.Fatal("can't prepare delete statement", err)
	}

	if _, err := stmt.Exec(id); err != nil {
		log.Fatal("can't execute delete statment", err)
	}
	done := Done{Message: "customer deleted"}
	c.JSON(http.StatusOK, done)
}
