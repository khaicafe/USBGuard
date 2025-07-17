# Remove-AppxPackage for Current User
# Created by Nguyen Tuan
# Website:  www.blogthuthuatwin10.com

If (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")) {
	Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`" $args" -Verb RunAs
	Exit
}

Function Main-menu()
{
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
            $IDs=Read-Host -Prompt "For remove each app please select ID and press enter"
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
			$AppName=$apps[$ID].packagefullname

			Remove-AppxPackage -Package $AppName
		    }
		    else
		    {
			echo ""
            Write-warning -Message "wrong ID"
            echo ""
            pause
		    }
        }
}
Main-menu