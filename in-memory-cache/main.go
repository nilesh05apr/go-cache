package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jellydator/ttlcache/v2"
)

type Todo struct {	
	UserId int `json:"userId"`
	Id int `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}

//I am trying to cache the response of the API call. I am using the ttlcache package. I am able to cache the response but I am not able to get the cached response. I am getting the response from the API call. I am not sure what I am doing wrong. I am new to Go and Fiber. Any help would be appreciated.
var cache ttlcache.SimpleCache = ttlcache.NewCache()

func main() {
	app := fiber.New()
	var notFound = ttlcache.ErrNotFound

	cache.SetTTL(10*time.Second)
	
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		val, err := cache.Get(id)
		if err != notFound {	
			var todo Todo
			response, err := http.Get("https://jsonplaceholder.typicode.com/todos/" + id)
			if err != nil {
				return c.JSON(fiber.Map{"error": err})
			}
			defer response.Body.Close()
			data, err := io.ReadAll(response.Body)
			if err != nil {
				return c.JSON(fiber.Map{"error": err})
			}
			json.Unmarshal(data, &todo)
			cache.Set(id, todo)
			return c.JSON(fiber.Map{"data": todo})
		}
		return c.JSON(fiber.Map{"cache-data": val})
	})

	app.Listen(":3000")
}