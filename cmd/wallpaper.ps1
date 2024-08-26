#Requires -Version 5.0

# Define script parameters
param (
    [string]$SourceFolderPath = "G:\StarGate\apod_images\favorites",  # Folder path containing the lockscreen images
    [int]$IntervalInSeconds = 30  # Slideshow interval in seconds
)

# Import the PSADT module
Import-Module "$PSScriptRoot\Toolkit\AppDeployToolkitMain.ps1" -Force

# Define the main installation function
function Install-MyLockscreenSlideshow {
    # Create a temporary folder to copy the lockscreen images
    $tempFolder = Join-Path -Path $env:TEMP -ChildPath "LockscreenImages"
    New-Item -ItemType Directory -Path $tempFolder | Out-Null

    try {
        # Copy the lockscreen images to the temporary folder
        Copy-Item -Path $SourceFolderPath -Destination $tempFolder -Recurse -Force

        # Configure the Lockscreen slideshow
        Set-ItemProperty -Path "HKCU:\Control Panel\Personalization\Desktop Slideshow" -Name "ImagesRootPath" -Value $tempFolder
        Set-ItemProperty -Path "HKCU:\Control Panel\Personalization\Desktop Slideshow" -Name "Interval" -Value $IntervalInSeconds

        # Refresh the Lockscreen settings
        $result = Invoke-Expression -Command 'RUNDLL32.EXE USER32.DLL,UpdatePerUserSystemParameters ,1 ,True'
        
        if ($result -eq 0) {
            Write-Host "Lockscreen slideshow configuration applied successfully."
        } else {
            Write-Host "Lockscreen slideshow configuration failed."
        }
    }
    finally {
        # Cleanup the temporary folder
        Remove-Item -Path $tempFolder -Recurse -Force
    }
}

# Call the main installation function
Install-MyLockscreenSlideshow

# Execute the PSADT functions
If (-not (Test-RunningInAppVirt)) {
    # Execute the main installation function
    Execute-Installation
}
