// secure_app.go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	driveLetter = "E"
	sector      = 2048
	aesKey      = "mysecretkey123456" // Phải giống với writer
	blockSize   = 512
)

func getUUID() (string, error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`(Get-Volume -DriveLetter %s).UniqueId`, driveLetter))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("cannot get UUID: %s", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func decrypt(encrypted []byte, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	iv := encrypted[:aes.BlockSize]
	encryptedText := encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encryptedText, encryptedText)

	return string(encryptedText), nil
}

func readFromSector() ([]byte, error) {
	path := fmt.Sprintf(`\\.\%s:`, driveLetter)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open disk: %s", err)
	}
	defer file.Close()

	offset := int64(sector * blockSize)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, fmt.Errorf("seek failed: %s", err)
	}

	buf := make([]byte, blockSize)
	_, err = file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("read failed: %s", err)
	}

	return buf, nil
}

func main() {
	expectedUUID, err := getUUID()
	if err != nil {
		fmt.Println("❌ Không thể lấy UUID:", err)
		return
	}

	data, err := readFromSector()
	if err != nil {
		fmt.Println("❌ Không thể đọc sector:", err)
		return
	}

	decrypted, err := decrypt(data, aesKey)
	if err != nil {
		fmt.Println("❌ Không giải mã được dữ liệu:", err)
		return
	}

	if decrypted != expectedUUID {
		fmt.Println("🚫 Không đúng USB – từ chối chạy.")
		return
	}

	fmt.Println("✅ USB hợp lệ. Phần mềm bắt đầu...")
	// 👇 Chạy logic chính của phần mềm bạn tại đây
}
