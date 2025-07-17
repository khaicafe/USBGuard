# Disable-WindowsOptionalFeature.ps1
# Created by Nguyen Tuan
# Website:  www.blogthuthuatwin10.com

If (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")) {
	Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`" $args" -Verb RunAs
	Exit
}

Function Main-menu()
{
	$index=1
	$Features=Get-WindowsOptionalFeature -Online | ? state -eq 'enabled' | select featurename | sort -Descending
	#return entire listing of features
	    Write-Host "ID`t Feature Name"
    	echo ""
	foreach ($Feature in $Features)
	{
		Write-Host " $index`t $($Feature.FeatureName)"
		$index++
	}
    
    Do
    {
        echo ""
        $IDs=Read-Host -Prompt "For disable each feature please select ID and press enter"
    }
    While($IDs -eq "")
    
	#check whether input values are correct
	try
	{	
		[int[]]$IDs=$IDs -split ","
	}
	catch
	{
		 Write-warning -Message "wrong ID."
	}

	foreach ($ID in $IDs)
	{
		#check id is in the range
		if ($ID -ge 1 -and $ID -le $Features.count)
		{
			$ID--
			#Disable each feature
			$FeatureName=$Features[$ID].FeatureName

			 Disable-WindowsOptionalFeature -Online -FeatureName $FeatureName -NoRestart  
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
