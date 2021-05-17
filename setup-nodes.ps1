$network = "swan"

$nodes = @( 
    "51da.uk", 
    "51db.uk",
    "51dc.uk",
    "51dd.uk",
    "51de.uk"
)

$expiryDate = ((Get-date).AddDays(90)).ToString("yyyy-MM-dd")

Write-Output "Network: $($network)"
Write-Output "Expiry: $($expiryDate)"
$ok = Read-Host "Ok? [y/N]"

if ($ok -ne "y" -and $ok -ne "Y") {
    Break
}

## Set-up SWIFT access nodes as OWID creators
$nodes | ForEach-Object {
    Write-Output "51Degrees : $($_)" 
    $Response = Invoke-WebRequest -URI "http://$($_)/owid/register?name=51Degrees"
    Write-Output $Response.StatusCode
}


## Set-up SWAN participant as OWID creators
$dir = dir www | ?{$_.PSISContainer}
foreach ($d in $dir){
    Get-ChildItem www\$d -Filter config.json | 
        ForEach-Object {
            $c = Get-Content www\$d\config.json | ConvertFrom-Json
            Write-Output "$($c.Name) : $($d.Name)" 
            $Response = Invoke-WebRequest -URI "http://$($d.Name)/owid/register?name=$($c.Name)"
            Write-Output $Response.StatusCode
        }
}

## Set-up SWIFT access nodes
$nodes | ForEach-Object {
    Write-Output "51degrees : $($_)" 
    $Response = Invoke-WebRequest -URI "http://$($_)/swift/register?network=$($network)&expires=$($expiryDate)&role=0"
    Write-Output $Response.StatusCode
}

## Set-up SWIFT storage Nodes
$nodes | ForEach-Object {
    For ($i = 1; $i -le 30; $i++) {
        Write-Output "51degrees : $($i).$($_)" 
        $Response = Invoke-WebRequest -URI "http://$($i).$($_)/swift/register?network=$($network)&expires=$($expiryDate)&role=1"
        Write-Output $Response.StatusCode
    }
}