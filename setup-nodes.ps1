$nodes = @{
    "swan" = @( 
        "51da.uk", 
        "51db.uk",
        "51dc.uk",
        "51dd.uk",
        "51de.uk");
    "cisne" = @(
        "an.cisne-demo.es")
}

# All nodes will be valid from the current date.
$startDate = (Get-date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm")

# Set the nodes to expire in 30 years time.
$expiryDate = ((Get-date).ToUniversalTime().AddYears(30)).ToString("yyyy-MM-dd")

Write-Output "Network: $($network)"
Write-Output "Starts: $($startDate)" 
Write-Output "Expiry: $($expiryDate)"
$ok = Read-Host "Ok? [y/N]"

if ($ok -ne "y" -and $ok -ne "Y") {
    Break
}

# Set-up SWIFT access nodes as OWID creators
$nodes.Values | ForEach-Object {
    $_ | ForEach-Object {
        Write-Output "51Degrees : $($_)" 
        $Response = Invoke-WebRequest -URI "http://$($_)/owid/register?name=51Degrees"
        Write-Output $Response.StatusCode
    }
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
$nodes.Keys | ForEach-Object {
    $network = $_
    Write-Output "51degrees : $($network)" 
    $nodes[$network] | ForEach-Object {
        $Response = Invoke-WebRequest -URI "http://$($_)/swift/register?network=$($network)&starts=$($startDate)&expires=$($expiryDate)&role=0"
        Write-Output "$($_): $($Response.StatusCode)"
    }
}

## Set-up SWAN SWIFT storage Nodes
$network = "swan"
$nodes[$network] | ForEach-Object {
    For ($i = 1; $i -le 30; $i++) {
        $domain = "$($i).$($_)"
        Write-Output "51degrees : $($domain)" 
        $Response = Invoke-WebRequest -URI "http://$($domain)/swift/register?network=$($network)&starts=$($startDate)&expires=$($expiryDate)&role=1"
        Write-Output "$($domain): $($Response.StatusCode)"
    }
}

## Set-up Cisne SWIFT storage Node
$network = "cisne"
$domain = "sn.cisne-demo.es"
$Response = Invoke-WebRequest -URI "http://sn.cisne-demo.es/swift/register?network=$($network)&starts=$($startDate)&expires=$($expiryDate)&role=1"
Write-Output "51degrees : $($domain)" 
$Response = Invoke-WebRequest -URI "http://$($domain)/swift/register?network=$($network)&starts=$($startDate)&expires=$($expiryDate)&role=1"
Write-Output "$($domain): $($Response.StatusCode)"

## Set-up SWIFT sharing node
Write-Output "51degrees : s.51da.uk" 
$Response = Invoke-WebRequest -URI "http://s.51da.uk/swift/register?network=$($network)&starts=$($startDate)&expires=$($expiryDate)&role=2"
Write-Output $Response.StatusCode