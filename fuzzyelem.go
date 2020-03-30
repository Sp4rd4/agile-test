package fuzzyelem

import "github.com/pkg/errors"

func Search(id, sourceFile, targetFile string) error {
	if id == "" {
		return errors.New("empty id")
	}
	if sourceFile == "" {
		return errors.New("source file path is missing")
	}
	if targetFile == "" {
		return errors.New("target file path is missing")
	}
	return nil
}
