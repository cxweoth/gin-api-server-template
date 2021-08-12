package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// A function to read json file and return map[string]interface without fixed struct
func ReadUnstructuredJsonFile(filePath string) (map[string]interface{}, error) {

	// Init result
	var result map[string]interface{}

	// Open jsonFile
	jsonFile, err := os.Open(filePath)

	if err != nil {
		return result, err
	}

	// Defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// Read to byte
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Do unmarshal and write to result
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil

}
