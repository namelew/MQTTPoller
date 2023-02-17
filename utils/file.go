package utils

import (
	"log"
	"os"
)

func FileExists(file string) bool{
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetJsonFromFile(file string) string{
	if !FileExists(file){
		return ""
	}

	data,err := os.ReadFile(file)

	if err != nil {
		log.Fatal(err.Error())
	}

	return string(data)
}