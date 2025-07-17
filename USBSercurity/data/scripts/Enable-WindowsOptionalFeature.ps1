# Enable-WindowsOptionalFeature
# Created by Nguyen Tuan
# Website:  www.blogthuthuatwin10.com

If (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]"Administrator")) {
	Start-Process powershell.exe "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`" $args" -Verb RunAs
	Exit
}

Function Main-menu()
{
	$index=1
	$Features=Get-WindowsOptionalFeature -Online | ? State -eq 'disabled' | Select FeatureName
	#return entire listing of features 
	    Write-Host "ID`t Feature Name"
    	Write-Host ""
	foreach ($Feature in $Features)
	{
		Write-Host " $index`t $($Feature.FeatureName)"
		$index++
	}
    
    Do
    {
        Write-Host ""
        $IDs=Read-Host -Prompt "For enable each feature please select ID and press enter"
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
		if ($ID -ge 1 -and $ID -le $Features.count)
		{
			$ID--
			#Enable each feature
			$FeatureName=$Features[$ID].FeatureName

			Enable-WindowsOptionalFeature -Online -FeatureName $FeatureName -All -NoRestart
            pause
            cls
            Main-menu
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
