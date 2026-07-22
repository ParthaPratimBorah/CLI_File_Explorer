package compare

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// calculates the SHA256 hash of one file
func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return "", fmt.Errorf( "could not open file %s: %w", filePath, err )
	}

	defer file.Close()

	hashCalculator := sha256.New()

	_, err = io.Copy(
		hashCalculator,
		file,
	)

	if err != nil {
		return "", fmt.Errorf( "could not calculate hash for %s: %w", filePath, err )
	}

	hashValue := hashCalculator.Sum(nil)

	return fmt.Sprintf( "%x", hashValue, ), nil
}