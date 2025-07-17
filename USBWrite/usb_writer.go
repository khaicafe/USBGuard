// usb_writer.go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

const (
	driveLetter = "E"                 // Ổ USB
	sector      = 2048                // Sector để ghi dữ liệu
	aesKey      = "mysecretkey123456" // 16 byte AES key (phải giống trong app chính)
	blockSize   = 512                 // Size mỗi sector
)

func getUUID() (string, error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`(Get-Volume -DriveLetter %s).UniqueId`, driveLetter))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("cannot get UUID: %s", err)
	}
	uuid := strings.TrimSpace(string(output))
	return uuid, nil
}

func encrypt(text, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	padded := make([]byte, aes.BlockSize+len(text))
	iv := padded[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(padded[aes.BlockSize:], []byte(text))
	return padded, nil
}

func writeToSector(encrypted []byte) error {
	path := fmt.Sprintf(`\\.\%s:`, driveLetter)
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open disk: %s", err)
	}
	defer file.Close()

	offset := int64(sector * blockSize)
	_, err = file.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("seek failed: %s", err)
	}

	// pad to block size
	buf := make([]byte, blockSize)
	copy(buf, encrypted)

	_, err = file.Write(buf)
	return err
}

func main() {
	uuid, err := getUUID()
	if err != nil {
		fmt.Println("❌ Lỗi lấy UUID:", err)
		return
	}
	fmt.Println("🔐 UUID:", uuid)

	encrypted, err := encrypt(uuid, aesKey)
	if err != nil {
		fmt.Println("❌ Lỗi mã hóa:", err)
		return
	}

	err = writeToSector(encrypted)
	if err != nil {
		fmt.Println("❌ Lỗi ghi sector:", err)
		return
	}

	fmt.Println("✅ Đã ghi dữ liệu mã hóa vào sector USB thành công!")
}
