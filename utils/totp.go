package utils

import (
	"github.com/pquerna/otp/totp"
)

func GenerationNewTOTP(issuer, accountName string, qrcodeWidth, qrcodeHeight int) (secret string, pngBase64String string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}
	// 查查这个image.Image怎么用
	image, err := key.Image(qrcodeWidth, qrcodeHeight)
	if err != nil {
		return "", "", err
	}
	base64String, err := ImageToPngBase64String(image)
	if err != nil {
		return "", "", err
	}

	return key.Secret(), base64String, err
}

func ValidateTOTP(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}
