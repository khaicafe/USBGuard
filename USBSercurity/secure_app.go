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
	"time"
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

// 🔐 Tăng cường bảo mật runtime mỗi 15 giây
func startSecurityMonitor() {
	go func() {
		for {

			if controller.IsDebugged() {
				fmt.Println("🛑 Phát hiện hack (runtime) – thoát.")
				os.Exit(1)
			}
			time.Sleep(15 * time.Second)
		}
	}()
}

func printHelp() {
	fmt.Println("🔒 USBGuard Secure App")
	fmt.Println("Sử dụng:")
	fmt.Println("  secure_app.exe boot           - Tối ưu hệ thống Windows")
	fmt.Println("  secure_app.exe fastApp        - Tăng tốc Windows")
	fmt.Println("  secure_app.exe help           - Hiển thị hướng dẫn")
}

func main() {

	fmt.Println(` 
                                                                                                                                       

888    d8P  888               d8b                    .d888          
888   d8P   888               Y8P                   d88P"           
888  d8P    888                                     888             
888d88K     88888b.   8888b.  888  .d8888b  8888b.  888888  .d88b.  
8888888b    888 "88b     "88b 888 d88P"        "88b 888    d8P  Y8b 
888  Y88b   888  888 .d888888 888 888      .d888888 888    88888888 
888   Y88b  888  888 888  888 888 Y88b.    888  888 888    Y8b.     
888    Y88b 888  888 "Y888888 888  "Y8888P "Y888888 888     "Y8888  
                                                                    
                                                                    
                                                                                      
`)

	startSecurityMonitor() // 🔐 Kiểm tra nền nâng cao mỗi 15s

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	if devBuild {
		// dev
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
	} else {
		// release
		// check lisen
		expectedUUID, err := getUUID()
		if err != nil {
			fmt.Println("❌ Không thể lấy UUID:", err)
			return
		}
		serial, err := getVolumeSerial()
		if err != nil {
			fmt.Println("❌ Không thể lấy serial:", err)
			return
		}

		key := expectedUUID + ":" + serial

		// Gọi các anti patch rải rác
		checkAssets(key)
		verifyStartup(key)

		rawData, err := readFromSector()
		if err != nil {
			fmt.Println("❌ Không thể đọc sector:", err)
			return
		}

		decrypted, err := decryptAESGCM(rawData, []byte(aesKey))
		if err != nil {
			fmt.Println("❌ Không giải mã được:", err)
			return
		}

		if decrypted != key {
			fmt.Println("🚫 Không đúng USB – từ chối chạy.")
			return
		}

		// ✅ Bắt đầu app thật
		fmt.Println("✅ USB hợp lệ. Chạy phần mềm...")

		// Chạy kiểm tra nền ẩn → anti patch runtime
		slowVerify(key)
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

}
