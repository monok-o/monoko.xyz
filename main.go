package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type Post struct {
	Title string
	Path  string
}

func getPosts() []Post {
	var res []Post

	postFiles, err := ioutil.ReadDir("./posts")
	if err != nil {
		log.Fatal(err)
	}

	for _, post := range postFiles {
		name := strings.TrimSuffix(post.Name(), ".md")
		res = append(res, Post{
			Title: strings.Replace(name, "-", " ", -1),
			Path:  name,
		})

		file, err := os.Create("./views/articles/" + name + ".html")
		if err != nil {
			log.Fatal(err)
		}

		file.Write(parsePost(post.Name()))
	}

	return res
}

func parsePost(name string) []byte {
	file, err := ioutil.ReadFile("./posts/" + name)
	if err != nil {
		log.Fatal(err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	return markdown.ToHTML(file, parser, nil)
}

func init() {
	log.Print("test")
	getPosts()
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	posts := getPosts()

	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{}, "layouts/main")
	})

	app.Get("/projects", func(c *fiber.Ctx) error {
		return c.Render("projects", fiber.Map{}, "layouts/main")
	})

	app.Get("/blog/:id?", func(c *fiber.Ctx) error {
		if c.Params("id") != "" {
			return c.Render("articles/"+c.Params("id"), fiber.Map{}, "layouts/main")
		} else {
			return c.Render("blog", fiber.Map{
				"Posts": posts,
			}, "layouts/main")
		}
	})

	log.Fatal(app.Listen(":3000"))
}
