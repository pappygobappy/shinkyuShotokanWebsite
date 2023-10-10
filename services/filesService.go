package services

import "os"

func mkdirAllLocal(path string) error{
	return os.MkdirAll(path, 0700)
}

func MkdirAll(path string) error {
	if os.Getenv("ENV") == "local" {
		return mkdirAllLocal(path)
	} else {
		return mkdirAllLocal(path)
	}
}
