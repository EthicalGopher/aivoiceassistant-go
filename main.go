
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store = session.New()

func Filepath(c *fiber.Ctx) string {
	sess, err := store.Get(c)
	if err != nil {
		return "default"
	}
	// Ensure session exists and is valid
	if sess.Fresh() {
		if err := sess.Save(); err != nil {
			return "default"
		}
	}
	return sess.ID()
}

func TTS(input, dir, sessionID string) (string, error) {
	filename := fmt.Sprintf("speech_%s", sessionID)
	fullPath := filepath.Join(dir, filename+".mp3")

	// Create temp directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("directory creation failed: %w", err)
	}

	// Delete the existing file for this session if it exists
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to remove existing file: %w", err)
	}

	speech := htgotts.Speech{
		Folder:   dir,
		Language: voices.English,
		Handler:  &handlers.Native{},
	}

	if _, err := speech.CreateSpeechFile(input, filename); err != nil {
		return "", fmt.Errorf("speech generation failed: %w", err)
	}

	return fullPath, nil
}

func cleanupOldFiles() {
	for {
		time.Sleep(1 * time.Minute)
		files, err := os.ReadDir("temp")
		if err != nil {
			continue
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) != ".mp3" {
				continue
			}

			fullPath := filepath.Join("temp", file.Name())
			info, err := os.Stat(fullPath)
			if err != nil {
				continue
			}

			if time.Since(info.ModTime()) > 5*time.Minute {
				os.Remove(fullPath)
			}
		}
	}
}

func main() {
	go cleanupOldFiles()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowMethods: "GET",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		input := c.Query("input")
		if input == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing 'input' query parameter")
		}

		filepath := Filepath(c)
		result, err := TTS(input, "temp", filepath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		data, err := os.ReadFile(result)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error reading file")
		}

		c.Set("Content-Type", "audio/mpeg")
		c.Set("Accept-Ranges", "bytes")
		c.Set("Access-Control-Allow-Origin", "*")

		return c.Send(data)
	})

	log.Fatal(app.Listen("0.0.0.0:8089"))
}
