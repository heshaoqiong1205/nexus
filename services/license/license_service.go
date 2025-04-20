package license

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"nexus/models"
)

type LicenseService struct {
	licenseModels models.ILicenseModels
}

func NewLicenseService() *LicenseService {
	return &LicenseService{
		licenseModels: &models.LicenseModels{},
	}
}

func NewLicenseServiceWithModels(licenseModels models.ILicenseModels) *LicenseService {
	return &LicenseService{
		licenseModels: licenseModels,
	}
}

func (service *LicenseService) Authentication(licenseID string, token string) error {
	license, err := service.licenseModels.Get(licenseID)
	if err != nil {
		return err
	}
	if !license.Status {
		return errors.New("license is not active")
	}
	h := hmac.New(sha256.New, []byte(license.Key))
	h.Write([]byte(license.ID))
	hashValue := h.Sum(nil)
	encodedHash := base64.StdEncoding.EncodeToString(hashValue)
	if encodedHash == token {
		return nil
	}
	return errors.New("invalid token")
}
