package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
	"net/http"
	"projectGo/core"
)

var database *sql.DB

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func indexHandler(c echo.Context) error {
	rows, err := database.Query("select * from first_project_go.articles")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var posts []*core.Post
	for rows.Next() {
		p := core.Post{}
		err := rows.Scan(&p.Id, &p.Title, &p.Content)
		if err != nil {
			fmt.Println(err)
			continue
		}
		posts = append(posts, &p)
	}

	return c.Render(http.StatusOK, "index", posts)
}

func writeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "write", nil)
}

func editHandler(c echo.Context) error {
	id := c.Param("id")

	q := fmt.Sprintln("SELECT * FROM first_project_go.articles WHERE id = ", id)
	rows, err := database.Query(q)
	if err != nil {
		return err
	}

	post := core.Post{}
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Title, &post.Content)
		if err != nil {
			fmt.Println(err)
		}
	}

	return c.Render(http.StatusOK, "write", post)
}

func savePostHandler(c echo.Context) error {
	id := c.FormValue("id")
	title := c.FormValue("title")
	content := c.FormValue("content")

	var q string
	if id != "" {
		q = fmt.Sprintln("UPDATE first_project_go.articles SET title = ", title, ", content = ", content, " WHERE id = ", id)
	} else {
		q = fmt.Sprintln("INSERT first_project_go.articles (title, content) VALUES (", title, ", ", content, ")")
	}

	_, err := database.Exec(q)
	if err != nil {
		fmt.Println(err)
	}

	return c.Redirect(http.StatusFound, "/")
}

func deleteHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.Redirect(http.StatusNotFound, "/")
	}

	q := fmt.Sprintln("DELETE FROM first_project_go.articles WHERE id = ", id)
	if _, err := database.Exec(q); err != nil {
		return c.Redirect(http.StatusNotFound, "/")
	}

	return c.Redirect(http.StatusFound, "/")
}

func main() {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	db, err := sql.Open("mysql", "root:@/first_project_go")

	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	e.GET("/", indexHandler)
	e.GET("/write", writeHandler)
	e.GET("/edit/:id", editHandler)
	e.GET("/delete/:id", deleteHandler)
	e.POST("/SavePost", savePostHandler)

	e.Logger.Fatal(e.Start(":1015"))
}
