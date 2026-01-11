package input

import (
	"bufio"
	"os"
)

func ReadTargets(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var targets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}
