# Add-AppxPackage for Current User
# Created by Nguyen Tuan
# # Website:  www.blogthuthuatwin10.com

If (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")) {
	Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`" $args" -Verb RunAs
	Exit
}
    $apps= @(
    "Microsoft.XboxIdentityProvider"
    "Microsoft.XboxSpeechToTextOverlay"
    "Microsoft.Wallet"
    "Microsoft.Messaging"
    "Microsoft.StorePurchaseApp"
    "Microsoft.DesktopAppInstaller"
    "Microsoft.XboxGameOverlay"
    "Microsoft.People"
    "Microsoft.MicrosoftStickyNotes"
    "Microsoft.BingNews"
    "Microsoft.BingWeather"
    "Microsoft.Office.Sway"
    "Microsoft.RemoteDesktop"
    "Microsoft.ZuneMusic"
    "Microsoft.WindowsCamera"
    "Microsoft.ZuneVideo"
    "Microsoft.WindowsStore"
    "Microsoft.3DBuilder"
    "Microsoft.Microsoft3DViewer"
    "Microsoft.Windows.Photos"
    "Microsoft.MSPaint"
    "Microsoft.Getstarted"
    "Microsoft.WindowsAlarms"
    "Microsoft.WindowsSoundRecorder"
    "Microsoft.WindowsCalculator"
    "microsoft.windowscommunicationsapps"
    "Microsoft.XboxApp"
    "Microsoft.Office.OneNote"
    "Microsoft.WindowsMaps"
    "Microsoft.SkypeApp"
)
    foreach ($app in $apps)
	{
        Get-AppxPackage -allusers $app | Foreach {Add-AppxPackage -register "$($_.InstallLocation)\appxmanifest.xml" -DisableDevelopmentMode}
    }