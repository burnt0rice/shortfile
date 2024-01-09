package controllers

import (
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go-shortfile/utils"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

const filesDirectory = "./files/"
const tempFilesDirectory = "./tempFiles/"

func UploadFile(c *fiber.Ctx) error {
	utils.Log.Infow("UploadFile: Request Start", "Header", c.GetReqHeaders())
	if _, err := os.Stat(filesDirectory); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(filesDirectory, os.ModePerm)
		if err != nil {
			log.Println(err)
			utils.Log.Error("UploadFile: " + err.Error())
		}
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		fmt.Println("No file given")
		utils.Log.Warn("UploadFile: No file given")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No file given",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		fmt.Println("Can't open file")
		utils.Log.Warn("UploadFile: Can't open file")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Can't open file",
		})
	}
	defer file.Close()

	uuid := uuid.NewString()
	t := time.Now()
	timestamp := t.Format("20060102150405") //YYYYMMDDhhmmss

	encodeId := base64.StdEncoding.EncodeToString([]byte(uuid[0:5] + "--" + timestamp + "--" + fileHeader.Filename))

	gzippedFile, err := os.Create(filesDirectory + encodeId + ".gz")
	defer gzippedFile.Close()

	gzipWriter := gzip.NewWriter(gzippedFile)
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, file)
	if err != nil {
		fmt.Println("Can't copy file to gzippedFile")
		utils.Log.Warn("UploadFile: Can't copy file to gzippedFile")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong :(",
		})
	}

	gzipWriter.Flush()

	utils.Log.Infow("UploadFile: Request End", "FileId", encodeId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "File successfully uploaded",
		"fileId":      encodeId,
		"downloadUrl": string(c.Request().URI().Host()) + "/api/d/" + encodeId,
	})
}

func GetFileInfo(c *fiber.Ctx) error {
	utils.Log.Infow("GetFileInfo: Request Start", "Header", c.GetReqHeaders())
	encodeId := c.Params("id")
	if encodeId == "" {
		fmt.Println("No id given")
		utils.Log.Warn("GetFileInfo: No id given")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No id given",
		})
	}

	uuidFileName, err := base64.StdEncoding.DecodeString(encodeId)
	if err != nil {
		fmt.Println("Invalid id given")
		utils.Log.Warn("GetFileInfo: Invalid id given")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid id given",
		})
	}

	fullFileName := strings.Split(string(uuidFileName), "--")
	fileUuid := fullFileName[0]
	fileTimestamp := fullFileName[1]
	fileName := fullFileName[2]

	utils.Log.Infow("GetFileInfo: Request End", "FileId", encodeId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File found",
		"file": fiber.Map{
			"uuid":      fileUuid,
			"timestamp": fileTimestamp,
			"name":      fileName,
		},
	})
}

func DownloadFile(c *fiber.Ctx) error {
	utils.Log.Infow("DownloadFile: Request Start", "Header", c.GetReqHeaders())
	if _, err := os.Stat(tempFilesDirectory); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(tempFilesDirectory, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	encodeId := c.Params("id")
	if encodeId == "" {
		fmt.Println("No id given")
		utils.Log.Warn("DownloadFile: No id given")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No id given",
		})
	}

	uuidFileName, err := base64.StdEncoding.DecodeString(encodeId)
	if err != nil {
		fmt.Println("Invalid id given")
		utils.Log.Warn("DownloadFile: Invalid id given")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid id given",
		})
	}

	fullFileName := strings.Split(string(uuidFileName), "--")
	fileName := fullFileName[2]

	gzippedFile, err := os.Open(filesDirectory + encodeId + ".gz")
	if err != nil {
		fmt.Println("Can't open gzippedFile")
		utils.Log.Error("DownloadFile: Can't open gzippedFile", "Error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong :(",
		})
	}
	defer gzippedFile.Close()

	gzipReader, err := gzip.NewReader(gzippedFile)
	defer gzipReader.Close()

	uncompressedFile, err := os.Create(tempFilesDirectory + fileName)
	if err != nil {
		fmt.Println("Can't create uncompressed file")
		utils.Log.Error("DownloadFile: Can't create uncompressed file", "Error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong :(",
		})
	}
	defer uncompressedFile.Close()
	defer os.Remove(uncompressedFile.Name())

	_, err = io.Copy(uncompressedFile, gzipReader)
	if err != nil {
		fmt.Println("Can't copy file to gzippedFile")
		utils.Log.Error("DownloadFile: Can't copy file to gzippedFile", "Error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong :(",
		})
	}

	utils.Log.Infow("DownloadFile: Request End", "FileId", encodeId)
	return c.Status(fiber.StatusOK).Download(uncompressedFile.Name(), fileName)
}
