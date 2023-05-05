package util

import (
	"os"
	"path"

	"github.com/42LoCo42/emo/shared"
	"github.com/pkg/errors"
)

func TokenFilePath() (string, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return "", errors.Wrap(err, "could not get config dir")
	}

	return path.Join(configPath, shared.TOKEN_FILE), nil

}

func SaveToken(token []byte) error {
	tokenFilePath, err := TokenFilePath()
	if err != nil {
		return errors.Wrap(err, "could not get token file path")
	}

	if err := os.WriteFile(
		tokenFilePath,
		token,
		0600,
	); err != nil {
		return errors.Wrap(err, "could not write token file")
	}

	return nil
}

func LoadToken() (token []byte, err error) {
	tokenFilePath, err := TokenFilePath()
	if err != nil {
		return nil, errors.Wrap(err, "could not get token file path")
	}

	token, err = os.ReadFile(tokenFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read token file")
	}

	return token, nil
}
