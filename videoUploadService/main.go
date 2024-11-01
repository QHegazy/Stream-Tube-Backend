package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type ProgressTracker struct {
	UploadedBytes int64
	TotalBytes    int64
}

var (
	progressTrackers = make(map[string]*ProgressTracker)
	mu               sync.Mutex
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Post("/upload", uploadFile)
	app.Get("/progress/:id", websocket.New(handleProgress))

	log.Fatal(app.Listen(":3500"))
}

func uploadFile(c *fiber.Ctx) error {
	var re = regexp.MustCompile(`.*`)

	allowedExtensions := map[string]bool{
		".mp4": true,
		".mkv": true,
		".flv": true,
		".avi": true,
		".mov": true,
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).SendString("Failed to get form data")
	}

	if len(form.File) != 1 {
		return c.Status(400).SendString("Only one file is allowed")
	}

	for key, files := range form.File {
		if re.MatchString(key) {
			fileHeader := files[0]
			ext := filepath.Ext(fileHeader.Filename)

			if _, ok := allowedExtensions[ext]; !ok {
				return c.Status(400).SendString("Invalid file type. Only .mp4, .mkv, and .flv are allowed")
			}

			file, err := fileHeader.Open()
			if err != nil {
				return c.Status(500).SendString("Failed to open file")
			}
			defer file.Close()

			uploadID := uuid.NewString()

			tracker := &ProgressTracker{
				TotalBytes: fileHeader.Size,
			}
			mu.Lock()
			progressTrackers[uploadID] = tracker
			mu.Unlock()

			go saveFile(uploadID, tracker, file)

			return c.JSON(fiber.Map{
				"id": uploadID,
			}) // Send the upload ID back to the client
		}
	}

	return c.Status(400).SendString("Failed to process the upload")
}

func saveFile(uploadID string, tracker *ProgressTracker, file io.Reader) {
	out, err := os.Create(uploadID)
	if err != nil {
		log.Println("Failed to create file:", err)
		mu.Lock()
		delete(progressTrackers, uploadID)
		mu.Unlock()
		return
	}
	defer out.Close()

	reader := io.TeeReader(file, tracker)
	_, err = io.Copy(out, reader)
	if err != nil {
		log.Println("Failed to save file:", err)
		mu.Lock()
		delete(progressTrackers, uploadID)
		mu.Unlock()
		return
	}

	mu.Lock()
	delete(progressTrackers, uploadID)
	mu.Unlock()
}

func handleProgress(c *websocket.Conn) {
	uploadID := c.Params("id")

	for {
		mu.Lock()
		tracker, exists := progressTrackers[uploadID]
		mu.Unlock()

		if !exists {
			c.WriteMessage(websocket.CloseMessage, []byte{})
			break
		}

		progress := float64(atomic.LoadInt64(&tracker.UploadedBytes)) / float64(tracker.TotalBytes) * 100
		message := fmt.Sprintf("%.2f", progress)

		if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	c.Close()
}

func (p *ProgressTracker) Write(b []byte) (int, error) {
	n := len(b)
	atomic.AddInt64(&p.UploadedBytes, int64(n))
	return n, nil
}