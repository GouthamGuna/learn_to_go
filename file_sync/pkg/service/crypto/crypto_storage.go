package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

type FileStorage struct {
	UploadDir  string
	EncryptKey []byte
}

func NewFileStorage(uploadDir string, encryptionKey []byte) *FileStorage {
	os.MkdirAll(uploadDir, 0755)
	return &FileStorage{
		UploadDir:  uploadDir,
		EncryptKey: encryptionKey,
	}
}

// SaveEncryptedFile saves and encrypts an uploaded file
func (fs *FileStorage) SaveEncryptedFile(file io.Reader, filename string) (string, error) {
	// Create a temporary file to store uploaded data
	tempPath := filepath.Join(fs.UploadDir, "temp_"+filename)
	outFile, err := os.Create(tempPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	defer os.Remove(tempPath) // Clean up temp file

	// Copy uploaded data to temp file
	_, err = io.Copy(outFile, file)
	if err != nil {
		return "", err
	}

	// Generate encrypted filename
	hash := sha256.Sum256([]byte(filename))
	encryptedFilename := hex.EncodeToString(hash[:]) + ".enc"

	// Encrypt the file
	encryptedPath := filepath.Join(fs.UploadDir, encryptedFilename)
	err = EncryptFile(tempPath, encryptedPath, fs.EncryptKey)
	if err != nil {
		return "", err
	}

	return encryptedFilename, nil
}

// GetDecryptedFile returns a decrypted file reader
func (fs *FileStorage) GetDecryptedFile(encryptedFilename string) (io.ReadCloser, error) {
	encryptedPath := filepath.Join(fs.UploadDir, encryptedFilename)

	// Create temporary file for decryption
	tempPath := filepath.Join(fs.UploadDir, "decrypt_temp_"+encryptedFilename)
	defer os.Remove(tempPath) // Clean up after use

	// Decrypt the file
	err := DecryptFile(encryptedPath, tempPath, fs.EncryptKey)
	if err != nil {
		return nil, err
	}

	// Open the decrypted file
	file, err := os.Open(tempPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetDecryptedFileBytes returns decrypted file content as bytes
func (fs *FileStorage) GetDecryptedFileBytes(encryptedFilename string) ([]byte, error) {
	encryptedPath := filepath.Join(fs.UploadDir, encryptedFilename)

	// Read encrypted file
	encryptedData, err := os.ReadFile(encryptedPath)
	if err != nil {
		return nil, err
	}

	// Decrypt in memory
	decryptedData, err := DecryptBytes(encryptedData, fs.EncryptKey)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

// FileExists checks if a file exists
func (fs *FileStorage) FileExists(filename string) bool {
	path := filepath.Join(fs.UploadDir, filename)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
