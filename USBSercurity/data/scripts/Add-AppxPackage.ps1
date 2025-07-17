# Add-AppxPackage for Current User
# Created by Nguyen Tuan
# Website:  www.blogthuthuatwin10.com

If (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")) {
	Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`" $args" -Verb RunAs
	Exit
}

Function Main-menu()
{
    $index=1
	$apps=Get-AppxPackage -PackageTypeFilter Bundle -AllUsers
	#return entire listing of applications 
	    Write-Host "ID`t App name"
        Write-Host ""
	foreach ($app in $apps)
	{
		Write-Host " $index`t $($app.name)"
		$index++
	}
    if ($apps)
    {
		$index++
        Write-Host ""
	}
    else
    {
        Write-Host "Apps not found"
        Write-Host ""
        pause
    }  
    Do
    {
        Write-Host ""
        $IDs=Read-Host -Prompt "For add each app please select ID and press enter"
    }
    While($IDs -eq "")
    
	#check whether input values are correct
	try
	{	
		[int[]]$IDs=$IDs -split ","
	}
	catch
	{ 
		 Write-Host "Error:" $_.Exception.Message
	}

	foreach ($ID in $IDs)
	{
		#check id is in the range
		if ($ID -ge 1 -and $ID -le $apps.count)
		{
			$ID--
			#Remove each app
			$AppName=$apps[$ID].name

			Get-AppxPackage -AllUsers $AppName | Foreach {Add-AppxPackage -Register "$($_.InstallLocation)\appxmanifest.xml" -DisableDevelopmentMode}
		}
		else
		{
			Write-Host ""
            Write-warning -Message "wrong ID"
            Write-Host ""
            pause
		}
	}
}
Main-menu
