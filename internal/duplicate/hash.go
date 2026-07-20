package duplicate

import (
	"fmt"
	"hash"
	"io"
	"os"
	"crypto/sha256"
	"hash/crc32"
)

//calculate the hash
func calculateHash(filePath string, algorithm string) (string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return "", fmt.Errorf("could not open file for hashing: %w", err)
	}

	defer file.Close()

	var hashCalculator hash.Hash

	switch algorithm {
	case "sha256":
		hashCalculator = sha256.New()
	case "crc32":
		hashCalculator = crc32.NewIEEE()
	default:
		return "", fmt.Errorf("unsupported hash algorithm: %s", algorithm)
	}
	_, err = io.Copy(hashCalculator, file)

	if err != nil {
		return "", fmt.Errorf("could not calculate hash for %s: %w", filePath, err)
	}

	hashValue := hashCalculator.Sum(nil)
	return fmt.Sprintf("%x", hashValue), nil
}