package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	engine := html.New("./", ".html")

	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		var input struct {
			Nama_gambar string
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		gambar, err := c.FormFile("gambar")
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		fmt.Printf("Nama file: %s \n", gambar.Filename)
		fmt.Printf("Ukuran file (kb): %d \n", gambar.Size/1024)
		fmt.Printf("Mime type : %s \n", gambar.Header.Get("Content-type"))

		splitDots := strings.Split(gambar.Filename, ".")
		ext := splitDots[len(splitDots)-1]
		fmt.Printf("Extension : %s \n", ext)

		namaFile := fmt.Sprintf("%s.%s", time.Now().Format("2006-01-02-15-04-05"), ext)
		fmt.Printf("Nama file baru : %s \n", namaFile)

		fileHeader, _ := gambar.Open()
		defer fileHeader.Close()

		imageConfig, _, err := image.DecodeConfig(fileHeader)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		width := imageConfig.Width
		height := imageConfig.Height
		fmt.Printf("Width : %d \n", width)
		fmt.Printf("Height : %d \n", height)

		folderUpload := filepath.Join(".", "uploads")

		if err := os.MkdirAll(folderUpload, 0770); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		if err := c.SaveFile(gambar, "./uploads/"+namaFile); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"title":       input.Nama_gambar,
			"nama_gambar": namaFile,
			"message":     "Gambar berhasil diupload",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
