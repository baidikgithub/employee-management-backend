package main

import (
	"fmt"
	_ "log"
	"net/http"
	"strconv"

	sqlx "github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

type Employee struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Address string `json:"address" db:"address"`
}

type Response struct {
	Message string
	Status  bool
}

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres password=Shohom@789 dbname=db_employee sslmode=disable")

	if err = db.Ping(); err != nil {
		fmt.Println(err)
	}

	respon := Response{
		Message: "Success executing the query",
		Status:  true,
	}

	responError := Response{
		Message: "Error executing the query",
		Status:  false,
	}

	e := echo.New()

	e.Use(middleware.CORS())

	e.GET("/users", func(c echo.Context) error {
		rows, _ := db.Queryx("select * from users")

		var users []Employee

		for rows.Next() {
			place := Employee{}
			rows.StructScan(&place)
			users = append(users, place)
		}

		return c.JSON(http.StatusOK, users)
	})

	e.GET("/users/:id", func(c echo.Context) error {

		id, _ := strconv.Atoi(c.Param("id"))

		user := Employee{}
		err = db.Get(&user, "SELECT * FROM users WHERE id = $1", id)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, responError)
		}
		return c.JSON(http.StatusOK, user)

	})

	e.POST("/users", func(c echo.Context) error {
		reqBody := Employee{}
		c.Bind(&reqBody)

		_, err = db.NamedExec("insert into users(name, phone, address) values (:name, :phone, :address)", reqBody)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responError)
		}

		return c.JSON(http.StatusOK, respon)
	})

	e.PUT("/users/update/:id", func(c echo.Context) error {
		reqBody := Employee{}

		if err := c.Bind(&reqBody); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		id, _ := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}
		reqBody.Id = id
		_, errQuery := db.NamedExec("update users SET name= :name, phone= :phone, address= :address WHERE id= :id", reqBody)
		if errQuery != nil {
			return c.JSON(http.StatusInternalServerError, responError)
		}

		return c.JSON(http.StatusOK, respon)
	})

	e.DELETE("/users/delete/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responError)
		}
		return c.JSON(http.StatusOK, respon)
	})

	e.Logger.Fatal(e.Start(":8087"))
}
