package main
import (
	"os"
)


func readConfig(path string) ([]byte, error)  {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	const READ_BUFFER = 100

	buffer := make([]byte, READ_BUFFER)
	var fileContents []byte
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			return nil, err
		}
		fileContents = append(fileContents, buffer[0:bytesRead]...)
		if bytesRead < 100 {
			break
		}
	}
	return fileContents, nil
}
