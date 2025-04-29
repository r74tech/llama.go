package common

import (
	"io/ioutil"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	ba, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ba, nil
}
