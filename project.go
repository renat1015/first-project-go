package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"projectGo/models"
)

var posts map[string]*models.Post

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	return t.templates.ExecuteTemplate(w, name, data)
}

func indexHandler(c echo.Context) error {
	fmt.Println(posts)

	return c.Render(http.StatusOK, "index", posts)
}

func writeHandler(c echo.Context) error {

	return c.Render(http.StatusOK, "write", nil)
}

func editHandler(c echo.Context) error {
	id := c.Param("id")
	post, found := posts[id]
	if !found {
		return c.Redirect(http.StatusNotFound, "/")
	}

	return c.Render(http.StatusOK, "write", post)
}

func savePostHandler(c echo.Context) error {
	id := c.FormValue("id")
	title := c.FormValue("title")
	content := c.FormValue("content")

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id := GenerateId()
		post := models.NewPost(id, title, content)
		posts[post.Id] = post
	}

	return c.Redirect(http.StatusFound, "/")
}

func deleteHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.Redirect(http.StatusNotFound, "/")
	}
	delete(posts, id)

	return c.Redirect(http.StatusFound, "/")
}

func main() {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	posts = make(map[string]*models.Post, 0)

	e.GET("/", indexHandler)
	e.GET("/write", writeHandler)
	e.GET("/edit/:id", editHandler)
	e.GET("/delete/:id", deleteHandler)
	e.POST("/SavePost", savePostHandler)

	e.Logger.Fatal(e.Start(":1015"))
}
