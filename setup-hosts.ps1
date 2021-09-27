$HostsSampleContent = Get-Content -Path .\hosts-sample
$updateHosts = $true

ForEach($line in $HostsSampleContent) {
    if($line -eq "") {continue}
    
    $sel = Select-String -Path "$Env:WINDIR\System32\drivers\etc\hosts" -Pattern "$($line)"
    if($null -ne $sel) {
        $updateHosts = $false
        Break
    }
}

if ($updateHosts) 
{
    Add-Content -Path $Env:WINDIR\System32\drivers\etc\hosts -Value $HostsSampleContent
} else {
    Write-Host "No hosts written as your system hosts file already contains some SWAN entries."
}