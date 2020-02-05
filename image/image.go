package image

import (
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"regexp"
)

// Get returns the type of Image
func CheckImageType(p string) (string, error) {
	mime, _, err := mimetype.DetectFile(p)
	if err != nil {
		return "", err
	}

	reg := regexp.MustCompile(`image\/([a-zA-Z0-9\-\_]+)`)
	subs := reg.FindStringSubmatch(mime)
	if subsLen := len(subs); subsLen != 0 {
		mime = subs[subsLen-1]
	} else {
		return "", errors.New("invalid image type")
	}

	return mime, nil
}
