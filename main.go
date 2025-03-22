package main

import (
	"fmt"
	"os"

	// "github.com/EthicalGopher/rag/tts"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	// htgotts "github.com/hegedustibor/htgo-tts"
	// "github.com/hegedustibor/htgo-tts/handlers"
	// "github.com/hegedustibor/htgo-tts/voices"
)

var store = session.New()
var path string
func Filepath(c *fiber.Ctx) string {
	sess, err := store.Get(c)
	if err != nil {
		fmt.Println("Session error:", err)
		return "default" // Fallback if session fails
	}
	filepath := sess.ID()
	if err := sess.Save(); err != nil {
		fmt.Println("Failed to save session:", err)
	}
	return filepath
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowMethods: "POST,GET,DELETE",
	}))
	app.Delete("/delete",func(c*fiber.Ctx)error{
		err:=os.RemoveAll(path)
		if err!=nil{
			fmt.Println(err)
			return c.SendStatus(fiber.ErrBadRequest.Code)
		}
		store.Reset()
		return c.SendString("Success fully deleted")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// Get input from query parameter
		input := c.Query("input")
		if input == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing 'input' query parameter")
		}
	
		os.RemoveAll(path)
		fmt.Println("This is the path "+path)
		filepathArg := Filepath(c)
		folderpath := "temp"
	
	
		// err:=os.Truncate(path,0)
		// if err!=nil{
		// 	fmt.Println(err)
		// }
		var result string
		
		func() {
			
			result, err := TTS(input, folderpath,filepathArg)
			if err != nil {
				fmt.Println("TTS error:", err)
			}
			fmt.Println(result)
			path = result
		}()

	

		

		return c.SendString(result)
	})

	app.Listen(":8089")
}

