package utils

import "os"

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

	data,_ := os.ReadFile(file)

	return string(data)
}