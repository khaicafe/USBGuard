# Remove-Installedsoftwares.ps1
# Created by Nguyen Tuan
# Website:  www.blogthuthuatwin10.com

If (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")) {
	Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`" $args" -Verb RunAs
	Exit
}

Function Main-menu()
{
    $index=1
    $items = @(
        "HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\*"
        "HKLM:\Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\*"
    )
	$softwares=Get-ItemProperty $items
	#return entire listing of softwares 
	    Write-Host "ID`t Installed software"
        echo ""
	foreach ($software in $softwares)
	{
		Write-Host " $index`t $($software.displayname)"
        $index++
    }
    if ($softwares)
    {
		$index++
        echo ""
	}
    else
    {
        Write-Host "Software not found"
        echo ""
        pause
        exit
    }
        Do
        {
		    echo ""
            $IDs=Read-Host -Prompt "For uninstall each software please select ID and press enter"
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
		    if ($ID -ge 1 -and $ID -le $softwares.count)
		    {
			    $ID--
			    #Remove each app
			    $software=$softwares[$ID].UninstallString
			    echo ""
                echo "Uninstalling software...."
                Start-Process cmd -ArgumentList "/c ""$software""" -NoNewWindow -Wait
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