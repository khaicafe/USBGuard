package controller

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"usbguard/embed_assets"
)

func isAdmin() bool {
	// Kiểm tra quyền admin bằng cách chạy lệnh net session
	cmd := exec.Command("net", "session")
	// cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	return err == nil
}

func copySetACL() {
	arch := os.Getenv("PROCESSOR_ARCHITECTURE")
	var srcPath string
	if strings.Contains(arch, "64") {
		srcPath = "data/amd64/SetACL.exe"
	} else {
		srcPath = "data/x86/SetACL.exe"
	}

	destPath := "C:\\Windows\\System32\\SetACL.exe"
	input, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Println("❌ Failed to read SetACL:", err)
		return
	}

	err = os.WriteFile(destPath, input, 0755)
	if err != nil {
		fmt.Println("❌ Failed to copy SetACL to System32:", err)
	} else {
		fmt.Println("✅ SetACL copied to System32")
	}
}

func grantRegistryPermission() {
	key := "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Advanced"
	subKey := key + "\\DelayedApps"

	commands := [][]string{
		{"-on", key, "-ot", "reg", "-actn", "setowner", "-ownr", "n:Administrators"},
		{"-on", key, "-ot", "reg", "-actn", "ace", "-ace", "n:Administrators;p:full"},
		{"-on", subKey, "-ot", "reg", "-actn", "setowner", "-ownr", "n:Administrators"},
		{"-on", subKey, "-ot", "reg", "-actn", "ace", "-ace", "n:Administrators;p:full"},
	}

	for _, args := range commands {
		cmd := exec.Command("SetACL.exe", args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ SetACL error (%v): %s\n", args, string(output))
		}
	}
}

func getWindowsBuildNumber() (int, error) {
	out, err := exec.Command("reg", "query", `HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion`, "/v", "CurrentBuildNumber").Output()
	if err != nil {
		return 0, err
	}

	output := string(out)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "CurrentBuildNumber") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				buildNum := parts[len(parts)-1]
				return strconv.Atoi(buildNum)
			}
		}
	}
	return 0, fmt.Errorf("build number not found")
}

func applyWin10Tweaks() {
	fmt.Println("[+] Applying Windows 10 tweaks...")
	copySetACL()
	grantRegistryPermission() // ✅ thêm dòng này trước khi gọi script

	scriptContent := `$parentKey = 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced'
$subKey = $parentKey + '\DelayedApps'
$acl = Get-Acl $parentKey
$acl.SetOwner([System.Security.Principal.NTAccount]'Administrators')
Set-Acl -Path $parentKey -AclObject $acl
$rule = New-Object System.Security.AccessControl.RegistryAccessRule('Administrators','FullControl','ContainerInherit,ObjectInherit','None','Allow')
$acl.AddAccessRule($rule)
Set-Acl -Path $parentKey -AclObject $acl
if (-not (Test-Path $subKey)) { New-Item -Path $subKey -Force | Out-Null }
$aclSub = Get-Acl $subKey
$aclSub.SetOwner([System.Security.Principal.NTAccount]'Administrators')
Set-Acl -Path $subKey -AclObject $aclSub
$ruleSub = New-Object System.Security.AccessControl.RegistryAccessRule('Administrators','FullControl','ContainerInherit,ObjectInherit','None','Allow')
$aclSub.AddAccessRule($ruleSub)
Set-Acl -Path $subKey -AclObject $aclSub
Set-ItemProperty -Path $subKey -Name BoxedIoPriority -Value 0 -Type DWord
Set-ItemProperty -Path $subKey -Name BoxedPagePriority -Value 1 -Type DWord
Set-ItemProperty -Path $subKey -Name BoxedPriorityClass -Value 16384 -Type DWord
Set-ItemProperty -Path $subKey -Name Delay_Sec -Value 0 -Type DWord
`

	tmpFile := filepath.Join(os.TempDir(), "usbguard_win10.ps1")
	err := os.WriteFile(tmpFile, []byte(scriptContent), 0644)
	if err != nil {
		fmt.Println("❌ Cannot write PowerShell script:", err)
		return
	}

	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-NoProfile", "-File", tmpFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("⚠️ PowerShell execution error:", err)
	}
	fmt.Println("📄 PowerShell Output:\n", string(output))

	// Optional: Delete temp file
	os.Remove(tmpFile)

	// Step 2: TrustedInstaller NSudo call
	exeDir, _ := os.Executable()
	tools, err := embed_assets.ExtractAllAssetsTo(filepath.Dir(exeDir))
	if err != nil {
		log.Fatal("Extract failed:", err)
	}

	nsudo := tools["NSudoLC.exe"]
	fmt.Println("✔️ NSudo path:", nsudo)

	nsudoCmd := exec.Command(nsudo, "-U:T", "-P:E", "-ShowWindowMode:Hide", "PowerShell.exe", "-WindowStyle", "Hidden", "-Command",
		"if (-Not (Test-Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Advanced\\DelayedApps')) { New-Item -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Advanced\\DelayedApps' -Force | Out-Null }; "+
			"Set-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Advanced\\DelayedApps' -Name BoxedIoPriority -Value 0 -Type DWord; "+
			"Set-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Advanced\\DelayedApps' -Name BoxedPagePriority -Value 1 -Type DWord; "+
			"Set-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Advanced\\DelayedApps' -Name BoxedPriorityClass -Value 16384 -Type DWord; "+
			"Set-ItemProperty -Path 'HKLM:\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Advanced\\DelayedApps' -Name Delay_Sec -Value 0 -Type DWord;",
	)
	// nsudoCmd.Run()
	outputSudo, err := nsudoCmd.CombinedOutput()
	if err != nil {
		fmt.Println("⚠️ PowerShell error:", err)
	}
	fmt.Println("📄 Output:\n", string(outputSudo))
	os.Remove(nsudo)

	applyCommonTweaks()
}

func applyWin11Tweaks() {
	fmt.Println("[+] Applying Windows 11 tweaks...")

	script := `$RegPath = 'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced';
$acl = Get-Acl $RegPath;
$acl.SetOwner([System.Security.Principal.NTAccount]'Administrators');
Set-Acl -Path $RegPath -AclObject $acl;
if (-not (Test-Path $RegPath'\\DelayedApps')) { New-Item -Path $RegPath'\\DelayedApps' -Force | Out-Null };
Set-ItemProperty -Path $RegPath'\\DelayedApps' -Name BoxedIoPriority -Value 0 -Type DWord;
Set-ItemProperty -Path $RegPath'\\DelayedApps' -Name BoxedPagePriority -Value 1 -Type DWord;
Set-ItemProperty -Path $RegPath'\\DelayedApps' -Name BoxedPriorityClass -Value 16384 -Type DWord;
Set-ItemProperty -Path $RegPath'\\DelayedApps' -Name Delay_Sec -Value 0 -Type DWord;`

	// exec.Command("powershell", "-NoProfile", "-WindowStyle", "Hidden", "-Command", script).Run()
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("⚠️ PowerShell error:", err)
	}
	fmt.Println("📄 Output:\n", string(output))

	applyCommonTweaks()
}

func applyCommonTweaks() {
	fmt.Println("[Common Tweaks đang được áp dụng...]")
	tweaks := [][]string{
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\policies\\system", "/v", "DelayedDesktopSwitchTimeout", "/t", "REG_DWORD", "/d", "0", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\policies\\system", "/v", "PromptOnSecureDesktop", "/t", "REG_DWORD", "/d", "1", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\policies\\system", "/v", "FilterAdministratorToken", "/t", "REG_DWORD", "/d", "0", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Winlogon", "/v", "EnableFirstLogonAnimation", "/t", "REG_DWORD", "/d", "0", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Multimedia\\SystemProfile", "/v", "SystemResponsiveness", "/t", "REG_DWORD", "/d", "0", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Multimedia\\SystemProfile", "/v", "NetworkThrottlingIndex", "/t", "REG_DWORD", "/d", "4294967295", "/f"},
		{"reg", "add", "HKLM\\SYSTEM\\CurrentControlSet\\Control\\Session Manager\\Memory Management", "/v", "LargeSystemCache", "/t", "REG_DWORD", "/d", "0", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Multimedia\\SystemProfile\\Tasks\\Games", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Multimedia\\SystemProfile\\Tasks\\Games", "/v", "Affinity", "/t", "REG_DWORD", "/d", "0", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Multimedia\\SystemProfile\\Tasks\\Games", "/v", "Background Only", "/t", "REG_SZ", "/d", "False", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Multimedia\\SystemProfile\\Tasks\\Games", "/v", "Clock Rate", "/t", "REG_DWORD", "/d", "10000", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Multimedia\\SystemProfile\\Tasks\\Games", "/v", "GPU Priority", "/t", "REG_DWORD", "/d", "8", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Multimedia\\SystemProfile\\Tasks\\Games", "/v", "Priority", "/t", "REG_DWORD", "/d", "2", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Multimedia\\SystemProfile\\Tasks\\Games", "/v", "Scheduling Category", "/t", "REG_SZ", "/d", "High", "/f"},
		{"reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Multimedia\\SystemProfile\\Tasks\\Games", "/v", "SFIO Priority", "/t", "REG_SZ", "/d", "High", "/f"},
		{"reg", "add", "HKLM\\SYSTEM\\CurrentControlSet\\Control\\FileSystem", "/v", "NtfsDisable8dot3NameCreation", "/t", "REG_DWORD", "/d", "0", "/f"},
		{"reg", "add", "HKLM\\SYSTEM\\CurrentControlSet\\Control\\FileSystem", "/v", "NtfsAllowExtendedCharacterIn8dot3Name", "/t", "REG_DWORD", "/d", "1", "/f"},
		{"reg", "add", "HKLM\\SYSTEM\\CurrentControlSet\\Control\\FileSystem", "/v", "ConfigFileAllocSize", "/t", "REG_DWORD", "/d", "500", "/f"},
	}

	for _, cmdArgs := range tweaks {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		// cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		err := cmd.Run()
		if err != nil {
			log.Printf("⚠️ Lỗi chạy lệnh: %v\n", cmdArgs)
		}
	}
}

func RunApps() {
	fmt.Println("🔧 Bắt đầu tối ưu Windows...")

	if !isAdmin() {
		log.Println("❌ Cần chạy dưới quyền Administrator.")
		return
	}

	build, err := getWindowsBuildNumber()
	if err != nil {
		log.Println("❌ Không lấy được Windows Build:", err)
		return
	}
	fmt.Println("🧱 Windows Build:", build)

	if build < 22000 {
		fmt.Println("🪟 Áp dụng tối ưu cho Windows 10...")
		applyWin10Tweaks()
	} else {
		fmt.Println("🪟 Áp dụng tối ưu cho Windows 11...")
		applyWin11Tweaks()
	}

	applyCommonTweaks()
	fmt.Println("✅ Tối ưu hệ thống hoàn tất.")
}

func HandleOther(args []string) {
	fmt.Println("⚙️ Đang xử lý tham số khác:", args)
	// Thêm logic riêng nếu cần ở đây
}
