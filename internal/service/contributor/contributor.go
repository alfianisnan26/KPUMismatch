package contributor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"kawalrealcount/internal/data/model"
)

type Service interface {
	FetchContributionData(token string) (model.ContributorData, error)
}

type Param struct {
	Secret string
}

type service struct {
	secret string
}

func (s service) FetchContributionData(token string) (model.ContributorData, error) {
	// Decode the token from base64
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return model.ContributorData{}, err
	}

	// Decrypt the token using AES
	block, err := aes.NewCipher([]byte(s.secret))
	if err != nil {
		return model.ContributorData{}, err
	}

	if len(decodedToken) < aes.BlockSize {
		return model.ContributorData{}, fmt.Errorf("cipher text too short")
	}

	iv := decodedToken[:aes.BlockSize]
	decodedToken = decodedToken[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(decodedToken, decodedToken)

	// Unmarshal the decrypted data into model.ContributorData struct
	var contributionData model.ContributorData
	if err := json.Unmarshal(decodedToken, &contributionData); err != nil {
		return model.ContributorData{}, err
	}

	return contributionData, nil
}

func New(param Param) (Service, error) {
	if param.Secret == "" {
		return nil, errors.New("empty Secret, please manually set on every build")
	}
	return &service{
		secret: param.Secret,
	}, nil
}
