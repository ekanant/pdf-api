package main

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"pdf-api/logs"
	"time"

	"github.com/gofiber/fiber/v2"
)

var httpClient = http.Client{
	Timeout: 600 * time.Second,
}

func main() {

	appBasePath := os.Getenv("API_ROOT")

	app := fiber.New()
	router := app.Group(appBasePath)
	router.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"staus": "ok",
		})
	})
	router.Post("/*", func(c *fiber.Ctx) error {
		/*
			Match all route and send to gotenburg
		*/
		reqBody := c.Body()
		reqContentType := string(c.Request().Header.ContentType())
		reqURI := string(c.Request().RequestURI())

		//Send request to gotenberg
		gotenbergUrl := "http://gotenberg:3000" + reqURI
		response, responseHeaders, err := proxyToGotenberg(gotenbergUrl, reqContentType, reqBody)
		//Set with header from gotenberg response
		for k, v := range responseHeaders {
			c.Set(k, v)
		}
		if err != nil {
			return c.JSON(map[string]interface{}{
				"error": err,
			})
		}

		{
			//Compress pdf with ghostscript
			pdfFileName := ""
			for k, v := range responseHeaders {
				c.Set(k, v)

				if k == "Content-Disposition" {
					//Get pdf filename for create tmp file
					_, params, _ := mime.ParseMediaType(v)
					pdfFileName = params["filename"]
				}
			}

			if pdfFileName != "" {
				_, err := os.Stat("tmp")
				if os.IsNotExist(err) {
					//Create tmp folder if not exist
					err := os.Mkdir("tmp", 0700)
					if err != nil {
						logs.Error(fmt.Sprintf("create tmp folder error, pdf_name=%s, error=%s", pdfFileName, err))
					}
				}

				originalPdfFilePath := "tmp/" + pdfFileName
				compressPdfFilePath := "tmp/compress_" + pdfFileName

				//Write file for compress process
				err = os.WriteFile(originalPdfFilePath, response, 0644)
				if err != nil {
					logs.Error(fmt.Sprintf("save tmp_file err, pdf_name=%s, error=%s", pdfFileName, err))
				}

				//Compress pdf
				pdfDpi := 350
				out, err := exec.Command(`/bin/sh`, `-c`, fmt.Sprintf(`sh ./scripts/shrinkpdf.sh %s %s %d`, originalPdfFilePath, compressPdfFilePath, pdfDpi)).Output()
				if err != nil {
					logs.Error(fmt.Sprintf("compress pdf err, pdf_name=%s, out=%s, error=%s", pdfFileName, out, err))
				}
				defer func() {
					//Delete original file
					err = os.Remove(originalPdfFilePath)
					if err != nil {
						logs.Error(fmt.Sprintf("cannot remove tmp_file err, pdf_name=%s, error=%s", pdfFileName, err))
					}
				}()

				defer func() {
					//Delete compress file
					err = os.Remove(compressPdfFilePath)
					if err != nil {
						logs.Error(fmt.Sprintf("cannot remove tmp_file err, pdf_name=%s, error=%s", compressPdfFilePath, err))
					}
				}()

				//Read compress pdf byte
				return c.SendFile(compressPdfFilePath)

			}
		}

		return c.Send(response)
	})

	app.Listen(":3000")
}

func proxyToGotenberg(gotenbergUrl string, requestContentType string, requestBody []byte) ([]byte, map[string]string, error) {
	resp, err := httpClient.Post(gotenbergUrl, requestContentType, bytes.NewReader(requestBody))
	if err != nil {
		// handle error
		logs.Error(fmt.Sprintf("send to gotenburg error=%s", err))
	}
	defer resp.Body.Close()

	responseHeaders := map[string]string{}

	for k, v := range resp.Header {
		responseHeaders[k] = v[0]
	}

	body, err := io.ReadAll(resp.Body)
	return body, responseHeaders, err
}
