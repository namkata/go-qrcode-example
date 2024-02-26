package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"qr-code-generator/qrcode"
	"qr-code-generator/utils"
)

func HandleRequest(writer http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(10 << 20)
	var size, content string = request.FormValue("size"), request.FormValue("content")
	var codeData []byte

	writer.Header().Set("Content-Type", "application/json")

	if content == "" {
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(
			"Could not determine the desired QR code content.",
		)
		return
	}

	qrCodeSize, err := strconv.Atoi(size)
	if err != nil || size == "" {
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode("Could not determine the desired QR code size.")
		return
	}

	qrCode := &qrcode.SimpleQRCode{Content: content, Size: qrCodeSize}
	watermarkFile, _, err := request.FormFile("watermark")
	if err != nil && errors.Is(err, http.ErrMissingFile) {
		codeData, err = qrCode.Generate()
		if err != nil {
			writer.WriteHeader(400)
			json.NewEncoder(writer).Encode(
				fmt.Sprintf("Could not generate QR code. %v", err),
			)
			return
		}
		writer.Header().Add("Content-Type", "image/png")
		writer.Write(codeData)
		return
	}

	watermark, err := utils.UploadFile(watermarkFile)
	if err != nil {
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(
			fmt.Sprint("Could not upload the watermark image.", err),
		)
		return
	}

	contentType := http.DetectContentType(watermark)
	if err != nil {
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(
			fmt.Sprintf(
				"Provided watermark image is a %s not a PNG. %v.", err, contentType,
			),
		)
		return
	}

	codeData, err = qrCode.GenerateWithWatermark(watermark)
	if err != nil {
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(
			fmt.Sprintf(
				"Could not generate QR code with the watermark image. %v", err,
			),
		)
		return
	}

	writer.Header().Set("Content-Type", "image/png")
	writer.Write(codeData)
}
