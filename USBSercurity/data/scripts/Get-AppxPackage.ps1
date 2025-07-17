# Get-AppxPackage for Current User
# Created by Nguyen Tuan
# Website:  www.blogthuthuatwin10.com

If (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")) {
	Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`" $args" -Verb RunAs
	Exit
}

	$index=1
	$apps=Get-AppxPackage -PackageTypeFilter Bundle
	#return entire listing of applications 
	Write-Host "ID`t App name"
    echo ""
	foreach ($app in $apps)
	{
		Write-Host " $index`t $($app.name)"
        $index++
    }
    if ($apps)
    {
		$index++
        echo ""
        pause
	}
    else
    {
        Write-Host "Apps not found"
        echo ""
        pause
    }
