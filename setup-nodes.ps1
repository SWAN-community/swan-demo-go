$networks = @{
    "swan-decentralized" = @( 
        @{
            accessNode = "51da.uk";
            name = "Operator-A";
            secret = "yes";
            scramble = "yes";
            storageNodes = @( "1.51da.uk", "2.51da.uk", "3.51da.uk", 
                                "4.51da.uk", "4.51da.uk", "6.51da.uk",
                                "7.51da.uk", "8.51da.uk", "9.51da.uk",
                                "10.51da.uk", "11.51da.uk", "12.51da.uk",
                                "13.51da.uk", "14.51da.uk", "15.51da.uk" ) 
        },
        @{
            accessNode = "51db.uk";
            name = "Operator-B";
            secret = "yes";
            scramble = "yes";
            storageNodes = @( "1.51db.uk", "2.51db.uk", "3.51db.uk", 
                                "4.51db.uk", "4.51db.uk", "6.51db.uk",
                                "7.51db.uk", "8.51db.uk", "9.51db.uk",
                                "10.51db.uk", "11.51db.uk", "12.51db.uk",
                                "13.51db.uk", "14.51db.uk", "15.51db.uk" ) 
        },
        @{
            accessNode = "51dc.uk";
            name = "Operator-C";
            secret = "yes";
            scramble = "yes";
            storageNodes = @( "1.51dc.uk", "2.51dc.uk", "3.51dc.uk", 
                                "4.51dc.uk", "4.51dc.uk", "6.51dc.uk",
                                "7.51dc.uk", "8.51dc.uk", "9.51dc.uk",
                                "10.51dc.uk", "11.51dc.uk", "12.51dc.uk",
                                "13.51dc.uk", "14.51dc.uk", "15.51dc.uk" ) 
        },
        @{
            accessNode = "51dd.uk";
            name = "Operator-D";
            secret = "yes";
            scramble = "yes";
            storageNodes = @( "1.51dd.uk", "2.51dd.uk", "3.51dd.uk", 
                                "4.51dd.uk", "4.51dd.uk", "6.51dd.uk",
                                "7.51dd.uk", "8.51dd.uk", "9.51dd.uk",
                                "10.51dd.uk", "11.51dd.uk", "12.51dd.uk",
                                "13.51dd.uk", "14.51dd.uk", "15.51dd.uk" ) 
        },
        @{
            accessNode = "51de.uk";
            name = "Operator-E";
            secret = "yes";
            scramble = "yes";
            storageNodes = @( "1.51de.uk", "2.51de.uk", "3.51de.uk", 
                                "4.51de.uk", "4.51de.uk", "6.51de.uk",
                                "7.51de.uk", "8.51de.uk", "9.51de.uk",
                                "10.51de.uk", "11.51de.uk", "12.51de.uk",
                                "13.51de.uk", "14.51de.uk", "15.51de.uk" ) 
        }
    )
    "swan-operator-a" = @( 
        @{
            accessNode = "single.51da.uk";
            name = "Operator-A";
            secret = "no";
            scramble = "no";
            cookieDomain = "51d.uk";
            storageNodes = @( "operator-a.51d.uk") 
        });
    "swan-operator-b" = @( 
        @{
            accessNode = "single.51db.uk";
            name = "Operator-B";
            secret = "no";
            scramble = "no";
            cookieDomain = "51d.uk";
            storageNodes = @( "operator-b.51d.uk") 
        });
    "swan-operator-c" = @( 
        @{
            accessNode = "single.51dc.uk";
            name = "Operator-C";
            secret = "no";
            scramble = "no";
            cookieDomain = "51d.uk";
            storageNodes = @( "operator-c.51d.uk") 
        });
    "swan-operator-d" = @( 
        @{
            accessNode = "single.51dd.uk";
            name = "Operator-D";
            secret = "no";
            scramble = "no";
            cookieDomain = "51d.uk";
            storageNodes = @( "operator-d.51d.uk") 
        });
    "swan-operator-e" = @( 
        @{
            accessNode = "single.51de.uk";
            name = "Operator-E";
            secret = "no";
            scramble = "no";
            cookieDomain = "51d.uk";
            storageNodes = @( "operator-e.51d.uk") 
        });
    "cisne" = @(
        @{
            accessNode = "an.cisne-demo.es";
            name = "Operator-Cisne";
            secret = "no";
            scramble = "no";
            cookieDomain = "sn.cisne-demo.es";
            storageNodes = @( "sn.cisne-demo.es") 
        }
    )
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

# Set-up all the SWIFT networks from the n
$networks.GetEnumerator() | ForEach-Object {
    $network = $_.Key
    Write-Output "Network: $($network)"
    $_.Value | ForEach-Object {

        # Set-up the OWID details for the operator.
        $name = $_.name
        $domain = $_.accessNode
        Write-Output "OWID: $($name): $($domain)" 
        $url = "http://$($domain)/owid/register?name=$($name)"
        $Response = Invoke-WebRequest -URI $url
        Write-Output $Response.StatusCode

        # Set-up the SWIFT access node for the operator.
        $secret = $_.secret
        $scramble = $_.scramble
        $cookieDomain = $_.cookieDomain
        Write-Output "SWIFT Access: $($name): $($domain) $($cookieDomain)" 
        $url = "http://$($domain)/swift/register?network=$($network)&" +
            "starts=$($startDate)&expires=$($expiryDate)&role=0&" +
            "secret=$($secret)&scramble=$($scramble)"
        $Response = Invoke-WebRequest -URI $url
        Write-Output $Response.StatusCode

        # Set-up the SWIFT storage nodes for the operator.
        $_.storageNodes | ForEach-Object {
            Write-Output "SWIFT Storage: $($name): $($_) $($cookieDomain)" 
            $url = "http://$($_)/swift/register?network=$($network)&" +
                "starts=$($startDate)&expires=$($expiryDate)&role=1&" +
                "secret=$($secret)&scramble=$($scramble)&" +
                "cookieDomain=$($cookieDomain)"
            $Response = Invoke-WebRequest -URI $url
            Write-Output $Response.StatusCode
        }
    }
}

## Set-up SWAN participant as OWID creators
$dir = dir www | ?{$_.PSISContainer}
foreach ($d in $dir) {
    Get-ChildItem www\$d -Filter config.json | 
        ForEach-Object {
            $c = Get-Content www\$d\config.json | ConvertFrom-Json
            Write-Output "$($c.Name) : $($d.Name)" 
            $Response = Invoke-WebRequest -URI "http://$($d.Name)/owid/register?name=$($c.Name)"
            Write-Output $Response.StatusCode
        }
}