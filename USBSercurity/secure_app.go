// secure_app.go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"usbguard/controller"
)

const (
	driveLetter = "E"
	sector      = 2048
	blockSize   = 512
	aesKey      = "mySuperSecretKeyAES256!" // Phải dài 16/24/32 byte cho AES-128/192/256
)

// Get UUID of the USB
func getUUID() (string, error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`(Get-Volume -DriveLetter %s).UniqueId`, driveLetter))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Get Volume Serial (extra check)
func getVolumeSerial() (string, error) {
	cmd := exec.Command("cmd", "/C", "vol", driveLetter+":")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	parts := strings.Split(string(out), "Serial Number is ")
	if len(parts) < 2 {
		return "", fmt.Errorf("cannot parse volume serial")
	}
	return strings.TrimSpace(parts[1]), nil
}

// AES-GCM decrypt
func decryptAESGCM(ciphertext []byte, key []byte) (string, error) {
	if len(ciphertext) < 12+16 {
		return "", fmt.Errorf("invalid ciphertext")
	}
	nonce := ciphertext[:12]
	tagged := ciphertext[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plain, err := aesgcm.Open(nil, nonce, tagged, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
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

// Hàm giả tên "vô hại", thật ra là Anti-Patch
func checkAssets(secret string) {
	sum := sha256.Sum256([]byte(secret + "_res"))
	if sum[2]^sum[15] != 0xAA {
		panic("💥 Data corrupted!")
	}
}

func verifyStartup(secret string) {
	sum := sha256.Sum256([]byte("boot:" + secret))
	if sum[7]+sum[9] != 0xEF {
		panic("💣 Tampering detected")
	}
}

func slowVerify(secret string) {
	go func() {
		sum := sha256.Sum256([]byte(secret + "::live"))
		if sum[0]^sum[1]^sum[2] != 0x42 {
			panic("☠️ Runtime integrity fail")
		}
	}()
}

func extractNSudoIfNotExist() {
	// go:embed assets/NSudoLC.exe
	var nsudoBytes []byte
	targetPath := "NSudoLC.exe"
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		err := os.WriteFile(targetPath, nsudoBytes, 0755)
		if err != nil {
			fmt.Println("❌ Failed to write NSudoLC.exe:", err)
		} else {
			fmt.Println("✅ Extracted NSudoLC.exe")
		}
	}
}

func printHelp() {
	fmt.Println("🔒 USBGuard Secure App")
	fmt.Println("Sử dụng:")
	fmt.Println("  secure_app.exe boot           - Tối ưu hệ thống Windows")
	fmt.Println("  secure_app.exe fastApp        - Tăng tốc Windows")
	fmt.Println("  secure_app.exe help           - Hiển thị hướng dẫn")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}
	// expectedUUID, err := getUUID()
	// if err != nil {
	// 	fmt.Println("❌ Không thể lấy UUID:", err)
	// 	return
	// }
	// serial, err := getVolumeSerial()
	// if err != nil {
	// 	fmt.Println("❌ Không thể lấy serial:", err)
	// 	return
	// }

	// key := expectedUUID + ":" + serial

	// // Gọi các anti patch rải rác
	// checkAssets(key)
	// verifyStartup(key)

	// rawData, err := readFromSector()
	// if err != nil {
	// 	fmt.Println("❌ Không thể đọc sector:", err)
	// 	return
	// }

	// decrypted, err := decryptAESGCM(rawData, []byte(aesKey))
	// if err != nil {
	// 	fmt.Println("❌ Không giải mã được:", err)
	// 	return
	// }

	// if decrypted != key {
	// 	fmt.Println("🚫 Không đúng USB – từ chối chạy.")
	// 	return
	// }

	// // ✅ Bắt đầu app thật
	// fmt.Println("✅ USB hợp lệ. Chạy phần mềm...")

	// // Chạy kiểm tra nền ẩn → anti patch runtime
	// slowVerify(key)

	// code logic app

	switch os.Args[1] {
	case "boot":
		controller.RunApplication()
	case "fastApp":
		controller.RunApps()
	case "help", "-h", "--help":
		printHelp()
	default:
		controller.HandleOther(os.Args[1:])
	}

}
