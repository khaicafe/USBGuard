# Restore-Computer
# Created by Nguyen Tuan
# Website:  www.blogthuthuatwin10.com

If (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")) {
	Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`" $args" -Verb RunAs
	Exit
}

    Get-WMIObject win32_service | Where-Object {$_.Name -eq "SDRSVC"}
    Set-Service SDRSVC -startuptype "manual"
    Start-Service SDRSVC
    Enable-ComputerRestore -Drive "C:\"
    Restore-Computer -RestorePoint 1