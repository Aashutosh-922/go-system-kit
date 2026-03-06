package main

import (
	"log"
	"time"

	"github.com/Aashutosh-922/go-system-kit.git/minigin"
)

func main() {
	app := minigin.New()

	app.Use(func(c *minigin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("%s %s (%s)", c.Request.Method, c.Request.URL.Path, time.Since(start))
	})

	app.GET("/health", func(c *minigin.Context) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	app.GET("/users/:id", func(c *minigin.Context) {
		c.JSON(200, map[string]string{"user_id": c.Param("id")})
	})

	app.GET("/files/*path", func(c *minigin.Context) {
		c.String(200, "requested path: "+c.Param("path"))
	})

	log.Println("mini gin demo listening on :8080")
	log.Println("try: curl localhost:8080/health")
	log.Println("try: curl localhost:8080/users/42")
	log.Println("try: curl localhost:8080/files/a/b/c.txt")

	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
