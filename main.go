package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	"github.com/EthicalGopher/rag/tts"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store = session.New()

func Filepath(c *fiber.Ctx) string {
	sess, err := store.Get(c)
	if err != nil {
		fmt.Println("Session error:", err)
		return "default"
	}
	filepath := sess.ID()
	if err := sess.Save(); err != nil {
		fmt.Println("Failed to save session:", err)
	}
	return filepath
}

func StartServer() *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowMethods: "POST,GET,DELETE",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		input := c.Query("input")
		if input == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing 'input' query parameter")
		}
		folderpath := "temp"
		filepathArg := Filepath(c)

		result, err := tts.TTS(input, folderpath, filepathArg)
		if err != nil {
			fmt.Println("TTS error:", err)
		}
		fmt.Println(result)

		return c.SendFile(result)
	})

	go func() {
		if err := app.Listen(":8089"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	return app
}

func main() {
	for {
		app := StartServer()

		// Wait for 1 minute (adjust to 1 hour for production)
		time.Sleep(1 * time.Minute)

		// Gracefully shut down the server first
		fmt.Println("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}

		// Clean up the "temp" directory after shutdown
		fmt.Println("Cleaning up temp directory...")
		err := deleteTempWithRetry("temp", 3, 2*time.Second)
		if err != nil {
			fmt.Println("Failed to clean up temp directory:", err)
		}

		// Restart the program
		fmt.Println("Restarting...")
		restartProgram()
	}
}

// Delete the temp directory with retry logic
func deleteTempWithRetry(dir string, maxRetries int, delay time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = os.RemoveAll(dir)
		if err == nil {
			return nil
		}
		fmt.Printf("Retry %d: Failed to delete %s: %v\n", i+1, dir, err)
		time.Sleep(delay)
	}
	return fmt.Errorf("failed to delete %s after %d retries: %v", dir, maxRetries, err)
}

// Restart the program
func restartProgram() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	cmd := exec.Command(executable, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to restart program: %v", err)
	}
	os.Exit(0)
}

