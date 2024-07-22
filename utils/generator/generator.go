package generator

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"math/big"
	"path/filepath"
)

func GenerateRandomCode() (string, error) {
	const charset = "0123456789"
	codeLength := 4

	randomBytes := make([]byte, codeLength)

	for i := range randomBytes {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		randomBytes[i] = charset[randIndex.Int64()]
	}

	return string(randomBytes), nil
}

func GenerateRandomBytes() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateEmailTemplate(fileTemplate string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("talkspace-api/utils/helper/email/template/%s", fileTemplate))
	if err != nil {
		return "", errors.New("invalid template file")
	}

	emailTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var templateBuffer bytes.Buffer
	if err := emailTemplate.Execute(&templateBuffer, data); err != nil {
		return "", err
	}

	return templateBuffer.String(), nil
}
