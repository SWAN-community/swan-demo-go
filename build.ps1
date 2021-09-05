# This is needed because the windows zip process used by EB will not enable
# the executable to be run on linux.
# https://forums.aws.amazon.com/message.jspa?messageID=825738REM825738
# go.exe get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip

# Set the common build environment variables.
$Env:GOPATH=""
$Env:GOARCH="amd64"

# Set up the AWS zip file command for Windows if it does not exist.
$zipcmd = "bin\build-lambda-zip.exe"
if (!(Test-Path $zipcmd))
{
    $Env:GOOS="windows"
    Invoke-Expression "go.exe install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@v1.26.0"
}

# Set up Go for AWS Elastic Beanstalk (EB) build
$Env:GOOS="linux"

# Build the application
Invoke-Expression "go build -o ./application ./server.go"

# Get all the files in the www folder that form the content with some files
# that are not needed for a live site removed.
$www = Get-ChildItem -File -Path ./www -Recurse | `
    Where-Object { $_.DirectoryName -NotMatch "\\test\\" } | `
    Where-Object { $_.DirectoryName -NotMatch "\\integrationExamples\\" } | `
    Where-Object { $_.Extension -NotMatch "md" } | `
    Where-Object { $_.Extension -NotMatch "yml" } | `
    Where-Object { $_.Name -ne "LICENSE" } | `
    Resolve-Path -Relative | `
    ForEach-Object { $a = $_ -replace '"', '""'; "`"$a`"" }

# Create the zip command with all the files.
if (Test-Path "application")
{
    $command = "build-lambda-zip.exe -o aws-eb-swan-demo.zip application appsettings.json Procfile .ebextensions/.config .ebextensions/healthcheckurl.config .swan/owidcreators-production.json .swan/swiftnodes-production.json " + $www -join ' '
    $command = $command.Replace(".\", "").Replace("\", "/")

    # Create a zip file with the application and the settings file
    Invoke-Expression $command
}

# Build the docker container image.
docker build -t swan-community/swan-demo-go .
