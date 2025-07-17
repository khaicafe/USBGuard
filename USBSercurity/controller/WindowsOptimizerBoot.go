package controller

import (
	"fmt"
	"os/exec"
	"strings"
)

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Run()
}

func runCommandSilent(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Run()
}

func getNetworkGUID() (string, error) {
	cmd := exec.Command("getmac", "/v", "/fo", "list")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Transport Name") && strings.Contains(line, "{") {
			start := strings.Index(line, "{")
			end := strings.Index(line, "}")
			if start != -1 && end != -1 {
				return line[start+1 : end], nil
			}
		}
	}
	return "", fmt.Errorf("Không tìm thấy GUID")
}

func RunApplication() {
	fmt.Println("🚀 Bắt đầu tối ưu hệ thống Windows...")

	guid, err := getNetworkGUID()
	if err != nil {
		fmt.Println("❌ Không tìm được GUID mạng:", err)
		return
	}

	// 1. Disable Nagle Algorithm
	fmt.Println("[1] Disable Nagle...")
	regPath := `HKLM\SYSTEM\CurrentControlSet\Services\Tcpip\Parameters\Interfaces\{` + guid + `}`
	runCommand("reg", "add", regPath, "/v", "TcpAckFrequency", "/t", "REG_DWORD", "/d", "1", "/f")
	runCommand("reg", "add", regPath, "/v", "TCPNoDelay", "/t", "REG_DWORD", "/d", "1", "/f")

	// 2. Hibernate Cleanup
	fmt.Println("[2] Reset Hibernate...")
	runCommandSilent("powercfg", "-h", "off")
	runCommandSilent("powercfg", "-h", "on")

	fmt.Println("[3] Dọn...")

	// 4. Registry Performance Tweaks
	fmt.Println("[4] Tối ưu Registry...")
	runCommand("reg", "add", `HKLM\SYSTEM\CurrentControlSet\Control\WMI\Autologger\ReadyBoot`, "/v", "Start", "/t", "REG_DWORD", "/d", "0", "/f")
	runCommand("reg", "add", `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Windows`, "/v", "LoadAppInit_DLLs", "/t", "REG_DWORD", "/d", "0", "/f")
	runCommand("reg", "add", `HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Memory Management`, "/v", "DisablePagingExecutive", "/t", "REG_DWORD", "/d", "1", "/f")
	runCommand("bcdedit", "/timeout", "0")
	runCommand("reg", "add", `HKLM\SYSTEM\CurrentControlSet\Control\PnP`, "/v", "DeviceInitTimeout", "/t", "REG_DWORD", "/d", "3", "/f")

	// 5. Game Mode / DVR
	fmt.Println("[5] Tắt Game Mode / Xbox...")
	runCommand("reg", "add", `HKLM\SOFTWARE\Policies\Microsoft\Windows\GameDVR`, "/v", "AllowGameDVR", "/t", "REG_DWORD", "/d", "0", "/f")
	runCommand("reg", "add", `HKCU\System\GameConfigStore`, "/v", "GameDVR_Enabled", "/t", "REG_DWORD", "/d", "0", "/f")

	// 6. Cortana, Bing Search
	fmt.Println("[6] Vô hiệu Cortana & Bing Search...")
	runCommand("reg", "add", `HKLM\SOFTWARE\Policies\Microsoft\Windows\Windows Search`, "/v", "AllowCortana", "/t", "REG_DWORD", "/d", "0", "/f")
	runCommand("reg", "add", `HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Search`, "/v", "BingSearchEnabled", "/t", "REG_DWORD", "/d", "0", "/f")

	// 7. Windows Tips & Content Delivery
	fmt.Println("[7] Tắt nội dung quảng cáo & gợi ý...")
	runCommand("reg", "add", `HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\ContentDeliveryManager`, "/v", "SubscribedContent-338387Enabled", "/t", "REG_DWORD", "/d", "0", "/f")
	runCommand("reg", "add", `HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\ContentDeliveryManager`, "/v", "SubscribedContent-338388Enabled", "/t", "REG_DWORD", "/d", "0", "/f")

	// 8. Telemetry + Services
	fmt.Println("[8] Tắt Telemetry...")
	runCommand("reg", "add", `HKLM\SOFTWARE\Policies\Microsoft\Windows\DataCollection`, "/v", "AllowTelemetry", "/t", "REG_DWORD", "/d", "0", "/f")
	runCommandSilent("sc", "stop", "DiagTrack")
	runCommandSilent("sc", "config", "DiagTrack", "start=", "disabled")

	// 9. Scheduled Task Cleanups
	fmt.Println("[9] Tắt task nền không cần thiết...")
	tasks := []string{
		"Microsoft\\Windows\\Application Experience\\Microsoft Compatibility Appraiser",
		"Microsoft\\Windows\\Customer Experience Improvement Program\\Consolidator",
		"Microsoft\\Windows\\Customer Experience Improvement Program\\KernelCeipTask",
		"Microsoft\\Windows\\Customer Experience Improvement Program\\UsbCeip",
		"Microsoft\\Windows\\DiskDiagnostic\\Microsoft-Windows-DiskDiagnosticDataCollector",
	}
	for _, task := range tasks {
		runCommand("schtasks", "/Change", "/TN", task, "/Disable")
	}

	// 10. Firewall / Security Notifications
	fmt.Println("[10] Tắt cảnh báo Security Center...")
	runCommand("reg", "add", `HKLM\SOFTWARE\Microsoft\Security Center\Notifications`, "/v", "DisableFirewallNotifications", "/t", "REG_DWORD", "/d", "1", "/f")
	runCommand("reg", "add", `HKLM\SOFTWARE\Microsoft\Windows Defender Security Center\Notifications`, "/v", "DisableNotifications", "/t", "REG_DWORD", "/d", "1", "/f")
	runCommand("reg", "add", `HKLM\SOFTWARE\Microsoft\Windows Defender Security Center\Systray`, "/v", "HideSystray", "/t", "REG_DWORD", "/d", "1", "/f")

	// 11. Delay non-critical services
	fmt.Println("[11] Chuyển các dịch vụ không thiết yếu sang delayed...")
	services := []string{
		"wuauserv", "BITS", "WSearch", "Spooler", "bthserv", "WerSvc",
		"DiagTrack", "ShellHWDetection", "Themes", "TabletInputService", "Fax", "FontCache",
		"lmhosts", "TrkWks", "wercplsupport", "stisvc", "SysMain", "PcaSvc", "SSDPSRV", "WinRM", "RemoteRegistry",
	}
	for _, svc := range services {
		runCommandSilent("sc", "config", svc, "start=", "delayed-auto")
	}

	// 12. Cleanup Temp Files
	fmt.Println("[12] Xóa file tạm...")
	runCommandSilent("cmd", "/c", "del /f /s /q %temp%\\*")
	runCommandSilent("cmd", "/c", "del /f /s /q C:\\Windows\\Temp\\*")

	// 13. Remove bloatware apps
	fmt.Println("[13] Gỡ app rác...")
	runCommandSilent("powershell", "-Command", "Get-AppxPackage *xbox* | Remove-AppxPackage")
	runCommandSilent("powershell", "-Command", "Get-AppxPackage *onenote* | Remove-AppxPackage")
	runCommandSilent("powershell", "-Command", "Get-AppxPackage *skype* | Remove-AppxPackage")

	// 14. DISM Cleanups
	fmt.Println("[14] Dọn WinSxS...")
	// runCommand("dism", "/online", "/cleanup-image", "/startcomponentcleanup")
	// runCommand("dism", "/online", "/cleanup-image", "/restorehealth")
	// runCommand("dism", "/online", "/cleanup-image", "/startcomponentcleanup", "/resetbase")

	// 15. CompactOS - Nén hệ điều hành để tiết kiệm dung lượng
	fmt.Println("[15] Kích hoạt CompactOS (tùy chọn)...")
	// runCommand("compact", "/compactos:always")

	// 16. Restart Explorer
	fmt.Println("[16] Restart Explorer...")
	// runCommandSilent("taskkill", "/f", "/im", "explorer.exe")
	// runCommand("start", "explorer.exe")

	fmt.Println("✅ Tối ưu hoàn tất! Vui lòng khởi động lại máy để áp dụng toàn bộ.")
}
