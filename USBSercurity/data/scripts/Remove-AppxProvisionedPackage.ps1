# Remove-AppxProvisionedPackage.ps1
# Created by Nguyen Tuan
# Website:  www.blogthuthuatwin10.com

If (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")) {
	Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`" $args" -Verb RunAs
	Exit
}

Function Main-menu()
{
	$index=1
	$apps=Get-AppxProvisionedPackage -online
	#return entire listing of applications 
	Write-Host "ID`t PakageName"
    echo ""
	foreach ($app in $apps)
	{
		Write-Host " $index`t $($app.packagename)"
		$index++
	}
    if ($apps)
    {
		$index++
	}
    else
    {
        Write-Host "Apps not found"
        echo ""
        pause
        exit
    }
        Do
        {
            echo ""
            $IDs=Read-Host -Prompt "For remove each package please select ID and press enter"
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
			#Remove each package
			$AppName=$apps[$ID].packagename

			Remove-AppxProvisionedPackage -Online -Package $AppName 
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