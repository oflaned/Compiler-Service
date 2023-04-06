package lib

import "os"

func CreateTempCFile(code string) (*os.File, error) {
	file, err := os.CreateTemp("./Temp", "code.*.cpp")
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(code)
	if err != nil {
		return nil, err
	}

	err = file.Chmod(0755)
	if err != nil {
		return nil, err
	}

	err = file.Sync()
	if err != nil {
		return nil, err
	}

	return file, nil
}
