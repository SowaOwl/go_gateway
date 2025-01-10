package debug

import (
	"encoding/json"
	"log"
	"os"
)

func LogJson(object interface{}) error {
	jsonResponse, err := json.Marshal(object)
	if err != nil {
		return err
	}

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	if _, err := file.Write([]byte("\n")); err != nil {
		return err
	}

	if _, err := file.Write(jsonResponse); err != nil {
		return err
	}

	return nil
}
