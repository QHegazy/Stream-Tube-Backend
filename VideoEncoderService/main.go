package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VideoInfo struct {
	Bitrate1080 int     `json:"bitrate_1080"`
	Bitrate720  int     `json:"bitrate_720"`
	Bitrate480  int     `json:"bitrate_480"`
	Bitrate360  int     `json:"bitrate_360"`
	// Duration    uint32 `json:"duration"`
}

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024,
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse form data")
		}
	
		// Retrieve the file headers. Check if there's exactly one file.
		files := form.File["video"]
		if len(files) != 1 {
			return c.Status(fiber.StatusBadRequest).SendString("Please upload exactly one file.")
		}
		
		file := files[0]
        pattern := `\.(mp4|mkv|flv)$`
        matched, err := regexp.MatchString(pattern, file.Filename)
        if err != nil {
            return err
        }

        if !matched {
            return c.Status(fiber.StatusBadRequest).SendString("Invalid file type. Only MP4, MKV, and FLV files are allowed.")
        }
		uuid := uuid.New().String()
		ext := filepath.Ext(file.Filename)

		uploadPath := "./uploads/" + uuid + ext
		if err := c.SaveFile(file, uploadPath); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to save video file")
		}

		// Query video information
		videoInfo, err := getVideoInfo(uploadPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to get video info: %v", err))
		}

		// Generate a UUID for the upload
	
		outputDir := "./uploads/" + uuid
		os.MkdirAll(outputDir, os.ModePerm)

		// Handle transcoding in a separate goroutine
		go func(uploadPath, outputDir string, videoInfo VideoInfo) {
			cmd := exec.Command("./transcode.sh", uploadPath, outputDir,
				strconv.Itoa(videoInfo.Bitrate360),
				strconv.Itoa(videoInfo.Bitrate480),
				strconv.Itoa(videoInfo.Bitrate720),
				strconv.Itoa(videoInfo.Bitrate1080))
			stdoutFile, err := os.Create(filepath.Join(outputDir, "transcode.log"))
			if err != nil {
				log.Printf("Failed to create stdout log file: %v", err)
				return
			}
			defer stdoutFile.Close()

			stderrFile, err := os.Create(filepath.Join(outputDir, "transcode_err.log"))
			if err != nil {
				log.Printf("Failed to create stderr log file: %v", err)
				return
			}
			defer stderrFile.Close()

			cmd.Stdout = stdoutFile
			cmd.Stderr = stderrFile

			if err := cmd.Run(); err != nil {
				log.Printf("Failed to transcode video: %v", err)
				return
			}
		}(uploadPath, outputDir, videoInfo)

		// Return only the UUID
		return c.JSON(fiber.Map{
			"uuid": uuid,
		})
	})

	app.Get("/status/:uuid", func(c *fiber.Ctx) error {
		uuid := c.Params("uuid")
		logFile := filepath.Join("./uploads", uuid, "transcode.log")

		progress, err := getTranscodingProgress(logFile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to get progress")
		}

		var  totalDuration uint32 = 39 ; 

		percentage, err := calculatePercentage(progress, totalDuration)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to calculate percentage")
		}

		return c.JSON(fiber.Map{
			"progress":   progress,
			"percentage": percentage,
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func getVideoInfo(filePath string) (VideoInfo, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filePath)

	output, err := cmd.Output()
	if err != nil {
		return VideoInfo{}, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return VideoInfo{}, err
	}

	format, ok := result["format"].(map[string]interface{})
	if !ok {
		return VideoInfo{}, fmt.Errorf("format information not found")
	}

	// Extract bitrate
	bitrateStr, ok := format["bit_rate"].(string)
	if !ok {
		return VideoInfo{}, fmt.Errorf("bitrate information not found")
	}

	bitrate, err := strconv.Atoi(bitrateStr)
	if err != nil {
		return VideoInfo{}, fmt.Errorf("failed to parse bitrate: %v", err)
	}
	bitrateK := bitrate / 1000

	bitrate1080 := int(float64(bitrateK) * 0.45)
	bitrate720 := int(float64(bitrateK) * 0.29)
	bitrate480 := int(float64(bitrateK) * 0.1)
	bitrate360 := int(float64(bitrateK) * 0.08)

	// Extract duration
	// durationStr, ok := format["duration"].(string)
	// if !ok {
	// 	return VideoInfo{}, fmt.Errorf("duration information not found")
	// }

	// duration, err := strconv.ParseUint(durationStr, 10, 32)
	// if err != nil {
	// 	return VideoInfo{}, fmt.Errorf("failed to parse duration: %v", err)
	// }

	return VideoInfo{
		Bitrate1080: bitrate1080,
		Bitrate720:  bitrate720,
		Bitrate480:  bitrate480,
		Bitrate360:  bitrate360,
		// Duration:    uint32(duration),
	}, nil
}

func getTranscodingProgress(logFilePath string) (string, error) {
	file, err := os.Open(logFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lastProgress string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "out_time_ms=") {
			lastProgress = line
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return lastProgress, nil
}

func calculatePercentage(progress string, totalDuration uint32) (float64, error) {
	parts := strings.Split(progress, "=")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid progress line format")
	}

	currentProgress, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse progress: %v", err)
	}

	if totalDuration == 0 {
		return 0, fmt.Errorf("total duration is zero")
	}

	percentage := float64(currentProgress) / float64(totalDuration) * 100
	return percentage, nil
}
