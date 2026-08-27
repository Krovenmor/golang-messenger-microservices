package http

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
)

const (
	minHeaderSize = 1024

	minImageWidth  = 100
	minImageHeight = 100
)

func GetDecodedImage(r *http.Request, imgMaxWidth, imgMaxHeight int) (image.Image, string, int64, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return nil, "", -1, err
	}
	defer file.Close()

	if header.Size <= minHeaderSize {
		log.Printf("GetDecodedImage: header.Size <= minHeaderSize, %d <= %d", header.Size, minHeaderSize)
		return nil, "", -1, ErrBadSize
	}

	head := make([]byte, 512)
	if _, err := file.Read(head); err != nil {
		log.Printf("GetDecodedImage: trouble with file.Read(head), err: %q", err)
		return nil, "", -1, ErrBadFile
	}

	mimeType := http.DetectContentType(head)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		log.Printf("GetDecodedImage: not supported mimeType, type: %q", mimeType)
		return nil, "", -1, ErrBadFile
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, "", -1, err
	}

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("GetDecodedImage: trouble with image.DecodeConfig(file), err: %q", err)
		return nil, "", -1, ErrBadFile
	}

	if config.Width > imgMaxWidth || config.Height > imgMaxHeight {
		if config.Width > imgMaxWidth {
			log.Printf("GetDecodedImage: config.Width > imgMaxWidth, %d > %d", config.Width, imgMaxWidth)
		}
		if config.Height > imgMaxHeight {
			log.Printf("GetDecodedImage: config.Height > imgMaxHeight, %d > %d", config.Height, imgMaxHeight)
		}
		return nil, "", -1, ErrBadSize
	}
	if config.Width < minImageWidth || config.Height < minImageHeight {
		if config.Width < minImageWidth {
			log.Printf("GetDecodedImage: config.Width < minImageWidth, %d < %d", config.Width, minImageWidth)
		}
		if config.Height < minImageHeight {
			log.Printf("GetDecodedImage: config.Height < minImageHeight, %d > %d", config.Height, minImageHeight)
		}
		return nil, "", -1, ErrBadSize
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, "", -1, err
	}

	img, format, err := image.Decode(file)
	if err != nil {
		log.Printf("GetDecodedImage: trouble with image.Decode(file), err: %q", err)
		return nil, "", -1, ErrBadFile
	}

	exactSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		log.Printf("GetDecodedImage: trouble with file.Seek(0, io.SeekEnd), err: %q", err)
		return nil, "", -1, ErrInteranl
	}

	return img, format, exactSize, nil
}
