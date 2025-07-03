package utils

import (
	"os"

	"github.com/spf13/afero"
)

var AppFs afero.Fs = afero.NewOsFs()

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func CheckForConfig() (bool, error) {
	filename := "skelly.yaml"

	_, err := AppFs.Stat(filename)
	if err == nil {
		// file exists
		return true, nil
	}
	if os.IsNotExist(err) {
		// file does _not_ exist
		return false, nil
	}
	return false, err
}
