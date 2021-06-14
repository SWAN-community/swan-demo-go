# ![Secured Web Addressability Network](https://raw.githubusercontent.com/SWAN-community/swan/main/images/swan.128.pxls.100.dpi.png)
# Secured Web Addressability Network (SWAN) Demo

Demo of SWAN, SWIFT and OWID implemented in Go.

## Quick Start 

To get you up and running quickly on a local machine using JSON for storage.

### Prerequisites 

* Local Go version 1.15 or greater installation sufficient to run the Go command
line.

### Steps

1. Clone the repository and navigate to the demo root directory:
   ```sh
   git clone --recurse-submodules https://github.com/SWAN-community/swan-demo-go
   cd swan-demo-go
   ```

2. Get dependencies:
    * **Linux**: Run `./dep.sh` in a terminal.
    * **Windows**: Run `.\dep.ps1` in Powershell.

3. Build the swan-demo server - if using VSCode then skip this step:
    * **Linux**: Run `./build.sh` in a terminal.
    * **Windows**: Run `go build -o .\application.exe .\src\server.go`

4. Configure hosts file, when running locally the demo requires that a lot of 
   host names be redirected to localhost - the following scripts will append 
   the contents of `hosts-sample` to your hosts file. Please review these 
   scripts before executing or update your hosts file manually.
   * **Linux**: Run `sudo ./setup-hosts.sh` in a 
     terminal. 
   * **Windows**: Run `.\setup-hosts.ps1` in an elevated Powershell terminal.

5. Set environment variables:
   * If using Visual Studio Code, then a launch file is provided for convenience. 
     Rename `.vscode\launch.json.rename` to `.vscode\launch.json` 
   * OR, set the following environment variables:

    **Linux**
    ```sh
    export PORT=80
    export OWID_FILE="swan/creators.json"
    export SWIFT_NODES_FILE="swan/swiftnodes.json"
    ```

    **Windows - CMD**
    ```bat
    setx PORT=80
    setx OWID_FILE="swan/creators.json"
    setx SWIFT_NODES_FILE="swan/swiftnodes.json"
    ```

    **Windows - Powershell**
    ```powershell
    $Env:PORT=80
    $Env:OWID_FILE="swan/creators.json"
    $Env:SWIFT_NODES_FILE="swan/swiftnodes.json"
    ```

6. Run the Demo Server:
   * **VSCode** If using Visual Studio Code, then the `.vscode\launch.json.rename` 
     file contains all the necessary settings to run and debug the demo.
   * **Linux**: Run `./application appsettings.dev.json` in a terminal.
   * **Windows**: Run `.\application.exe .\appsettings.dev.json` in a Powershell 
     window.

7. Navigate to http://new-pork-limes.uk in your preferred browser.

# SWAN Concepts

The SWAN demo implements the concepts explained in 
[SWAN](https://github.com/SWAN-community/swan)

The deployment guide later on in this readme shows you how to set up example 
with multiple SWIFT nodes, publishers and marketers but every SWAN needs at 
least the following:

* 5x SWIFT Nodes
* 1x SWAN Node
* 1x Publisher
* 1x Marketer

An example of nodes and members with example domain names:
```
             +-----------+  +-----------+  +-----------+
Swift        |           |  |           |  |           |
Storage      | 1.51da.uk |  | 2.51da.uk |  | 3.51da.uk |  ... etc
Nodes        |           |  |           |  |           |
             +-----------+  +-----------+  +-----------+

             +---------+  +---------+
Swan Nodes & |         |  |         |
Swift Access | 51da.uk |  | 51db.uk | ... etc
Nodes        |         |  |         |
             +---------+  +---------+

             +-------------------+  +----------------+
             |                   |  |                |
Publisher    | new-pork-limes.uk |  | current-bun.uk | ... etc
             |                   |  |                |
             +-------------------+  +----------------+

             +--------------+  +---------------+  +----------------+
             |              |  |               |  |                |
Marketer     | cool-cars.uk |  | cool-bikes.uk |  | cool-creams.uk | ... etc
             |              |  |               |  |                |
             +--------------+  +---------------+  +----------------+
```

# Files

`Procfile` : needed by AWS Elastic Beanstalk to indicate the application 
executable for web services.

`build.ps1` : builds AWS or Azure packages on Windows ready for manual deployment.

`build.sh` : builds AWS or Azure packages on Linux ready for manual deployment.

`cert.ps1` : Creates self signed certificates to run the demo using the https 
scheme. You will need `openssl` available in your path environment variable.

`dep.ps1` : gets SWAN demo dependencies on Windows.

`Dockerfile.rename` : A sample dockerfile which can be used to build a Docker
container suitable for use with Azure App Services.

`appsettings.json` : application settings for production.

`appsettings.dev.json` : application settings for development.

`.ebextensions/.config.rename` : AWS Elastic Beanstalk .config template ready for
additional SSL certificates.

`hosts-sample` : A sample hosts file containing all the domains used in the demo

`setup-hosts.sh` : Appends the contents of `hosts-sample` to the system hosts 
file (Linux).

`setup-hosts.ps1` : Appends the contents of `hosts-sample` to the system hosts 
file (Windows).

`setup-nodes.ps1` : Registers creators and set's up Swift access and storage 
nodes (Windows).

`setup-nodes.sh` : Registers creators and set's up Swift access and storage 
nodes (Linux).

`.vscode/launch.json.rename` : template Visual Studio Code launch settings 
including place holders for the storage environment variable values.

Note: `.gitignore` will ignore `launch.json`, and `.ebextensions/.config` to 
limit the risk of commits containing access keys.

# Deployment

The demo currently supports the following environments:

* AWS Elastic Beanstalk
* Azure App Service
* Local Go SDK

And the following storage solutions:

* AWS Dynamo DB
* Azure Storage Tables
* Google Firebase
* Local JSON files.

### Get the code

This demo uses submodules, to clone the repository and the submodules at the 
same time, run:

```
git clone --recurse-submodules https://github.com/SWAN-community/swan-demo-go
```

## Local Installation

### Prerequisites 

* Local Go version 1.15 or greater installation sufficient to run the Go command
line.

* Specify a storage solution, see [Deployment](#deployment) for supported 
  storage solutions.

### Steps

* Set up a storage option.

  * For local storage, files containing configured creators and swift nodes are 
    available in the `.swan` directory. Set the following environment variables to
    use local storage:

    ```                
    OWID_FILE: ".swan/creators.json"
    SWIFT_NODES_FILE: ".swan/swiftnodes.json"
    ```

    The vscode `launch.json` file can also be used to set environment variables, 
    see `the .vscode\launch.json.rename` sample file.

  * OR, for AWS:
    * Set the ``AWS_ENABLED`` environment variable to true.
    * Specify region and credentials in either environment variables 
    or in `~/.aws/credentials` and `~/.aws/config` files. See 
    [Configuring AWS SDK for Go](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html) to set up your environment. Make sure 
    to set the `AWS_REGION`, `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` as 
    a minimum.

  * OR, for Azure Storage:
    
    **Linux**:
    ```bash
    export AZURE_STORAGE_ACCOUNT="<youraccountname>"
    export AZURE_STORAGE_ACCESS_KEY="<youraccountkey>"
    ```
    **Windows**:

    Cmd:
    ```bat
    setx AZURE_STORAGE_ACCOUNT "<youraccountname>"
    setx AZURE_STORAGE_ACCESS_KEY "<youraccountkey>"
    ```
    Powershell:
    ```powershell
    $Env:AZURE_STORAGE_ACCOUNT = "<youraccountname>"
    $Env:AZURE_STORAGE_ACCESS_KEY = "<youraccountkey>"
    ```

  * OR, for Firebase see [Google FireBase](#google-firebase)

* Having cloned the repository, configure the Go environment. If using Visual
Studio Code see `.vscode/launch.json.rename` for a template environment. Rename 
to remove the `rename` extension suffix. As this file will contain sensitive 
information it is exclude in .gitignore.

* Configure your hosts file to point URLs to localhost, see 
[Environments > Local](#local) section in this readme for platform specifics. 
The following host resolutions are used in the sample configuration:

  ```
  # domains from the www folder
  127.0.0.1	new-pork-limes.uk
  127.0.0.1	new-prebid-limes.uk
  127.0.0.1 biscuit-news.uk
  127.0.0.1	current-bun.uk
  127.0.0.1	cool-bikes.uk
  127.0.0.1	cool-cars.uk
  127.0.0.1	cool-creams.uk
  127.0.0.1	cmp.swan-demo.uk
  127.0.0.1	swan-demo.uk
  127.0.0.1	pop-up.swan-demo.uk
  127.0.0.1	badssp.swan-demo.uk
  127.0.0.1	bidswitch.swan-demo.uk
  127.0.0.1	centro.swan-demo.uk
  127.0.0.1	dataxu.swan-demo.uk
  127.0.0.1	liveintent.swan-demo.uk
  127.0.0.1	magnite.swan-demo.uk
  127.0.0.1	mediamath.swan-demo.uk
  127.0.0.1	oath.swan-demo.uk
  127.0.0.1	pubmatic.swan-demo.uk
  127.0.0.1 sirdata.swan-demo.uk
  127.0.0.1	smaato.swan-demo.uk
  127.0.0.1	thetradedesk.swan-demo.uk
  127.0.0.1	zeta.swan-demo.uk
  127.0.0.1	liveramp.swan-demo.uk
  127.0.0.1	quantcast.swan-demo.uk
  # swift access nodes
  127.0.0.1	51da.uk
  127.0.0.1	51db.uk
  127.0.0.1	51dc.uk
  127.0.0.1	51dd.uk
  127.0.0.1	51de.uk
  # swift storage nodes
  127.0.0.1	1.51da.uk
  127.0.0.1	2.51da.uk
  127.0.0.1	3.51da.uk
  127.0.0.1	4.51da.uk
  127.0.0.1	5.51da.uk
  # Add more if your wish...
  ```
  More storage nodes can be found in the `hosts-sample` file.

* To build on:
  * **Linux**: Run the `./build.sh` file. 
  * **Windows**: First run `./dep.ps1` and then run the `./build.ps1` files.

* Run the server:

  ```
  ./application appsettings.dev.json
  ```

* Set up the access nodes, participants and storage nodes. Either run the setup 
  scripts or follow the steps to set up the nodes and participants manually.

  _NB: If you are using local storage then you can skip this step as the local 
  storage files in the `.swan` directory already contains a sample configuration._

  **Linux**: While the demo server is running, run the `./setup-nodes.sh` script 
  in the Terminal.

  **Windows**: While the demo server is running, run the `.\setup-nodes.ps1` 
  script in a Powershell Terminal.

  **Manually**:

  * The SWAN access nodes will be used to sign all the outgoing Open Web IDs and
  also to capture people's preferences. Register this domain with the following
  URL entering any of the details requested. This will create a record in the 
  ``owidcreators`` table for the domain which will contain randomly generated 
  public and private signing keys.

    http://51da.uk/owid/register \
    http://51db.uk/owid/register \
    http://51dc.uk/owid/register \
    http://51dd.uk/owid/register \
    http://51de.uk/owid/register 


  * The other SWAN participant nodes will also need to be registered as Open Web
  ID processors.

    http://new-pork-limes.uk/owid/register \
    http://biscuit-news.uk/owid/register \
    http://current-bun.uk/owid/register \
    http://cool-bikes.uk/owid/register \
    http://cool-cars.uk/owid/register \
    http://cool-creams.uk/owid/register \
    http://cmp.swan-demo.uk/owid/register \
    http://swan-demo.uk/owid/register \
    http://pop-up.swan-demo.uk/owid/register \
    http://badssp.swan-demo.uk/owid/register \
    http://bidswitch.swan-demo.uk/owid/register \
    http://centro.swan-demo.uk/owid/register \
    http://dataxu.swan-demo.uk/owid/register \
    http://liveintent.swan-demo.uk/owid/register \
    http://magnite.swan-demo.uk/owid/register \
    http://mediamath.swan-demo.uk/owid/register \
    http://oath.swan-demo.uk/owid/register \
    http://pubmatic.swan-demo.uk/owid/register \
    http://sirdata.swan-demo.uk/owid/register \
    http://smaato.swan-demo.uk/owid/register \
    http://thetradedesk.swan-demo.uk/owid/register \
    http://zeta.swan-demo.uk/owid/register \
    http://liveramp.swan-demo.uk/owid/register \
    http://quantcast.swan-demo.uk/owid/register

  * For each of the access and storage nodes that will be used for the SWIFT 
  component of the demo register these using the following URL. Enter the network 
  as "swan" (no quotes). The records from these steps will be visible in the 
  ``swiftnodes`` and ``swiftsecrets`` tables or if using local storage the 
  ``swiftnodes.json`` file.

    **Access nodes**

    http://51da.uk/swift/register \
    http://51db.uk/swift/register \
    http://51dc.uk/swift/register \
    http://51dd.uk/swift/register \
    http://51de.uk/swift/register

    **Storage nodes**

    http://1.51da.uk/swift/register \
    http://2.51da.uk/swift/register \
    http://3.51da.uk/swift/register \
    http://4.51da.uk/swift/register \
    http://5.51da.uk/swift/register \
    etc... 

* Now browse to one of the publisher URLs, you will be prompted to set your 
preferences:

  http://new-pork-limes.uk

# Environments

## AWS Elastic Beanstalk - without docker

### Prerequisites

* Local Go version 1.15 or greater installation sufficient to run the Go command
line.

* Familiar with the concepts associated with 
[SWAN](https://github.com/SWAN-community/swan),
[SWIFT](https://github.com/SWAN-community/swift), and 
[OWID](https://github.com/SWAN-community/owid).

* AWS account with Elastic Beanstalk and DynamoDB administration privileges or
Azure account with Storage and App Services.

#### AWS

* Setup domains to use with the demo in AWS using 
[Route 53](https://console.aws.amazon.com/route53). Domains are needed for the 
following roles.

    * SWAN access domain. For example ``51da.uk``.

    * Each of the SWIFT nodes that will support SWAN. At least two domains are 
    needed. For the purposes of the demo they may be sub domains. For example
    ``51da.uk`` and ``51db.uk``.

    * At least two different domains for publishers. For example 
    ``new-pork-limes.uk`` and ``current-bun.uk``.

    * At least two different domains for marketers. For example 
    ``cool-bikes.uk`` and ``cool-creams.uk``.

* Setup SSL certificates in AWS using 
[Certificate Manager](https://console.aws.amazon.com/acm/) 
for each of the domains.

#### Windows

* Set powershell unrestricted execution policy to enable `build.ps1` to execute.
  Use the following command with administrator privileges.

  `Set-ExecutionPolicy Unrestricted`

### Steps

* Get all the dependencies needed by the Go application.

  `go get -d ./...`

  `dep.ps1` can also be used in Powershell to explicitly get the dependencies 
  for this Go application.

* Add the demo domain names as folder names to the www folder. For example; the
domain `domain.com` would appear as `www/domain.com`. Alter the `config.json` 
content in each folder to indicate the purpose of the domain. This aspect of 
the demo is changing and domain examples provided with the should be reviewed 
along with the demo source code and comments to understand how multiple domains
are supported within a single demo application.

* You may need to support multiple SSL certificates if the demo deployment 
should respond to five or more domains. 
See AWS
[documentation](https://aws.amazon.com/premiumsupport/knowledge-center/elastic-beanstalk-ssl-load-balancer/).
Locate the ARN for each of the certificates that should be used with the demo.

* Add SSL certificate ARNs to a copy of ``.ebextensions/.config``. For example 
if you have the SSL ARNs A, B and C your .config file would contain the
following entries.

  ```
  option_settings:
    aws:elbv2:listener:443:
      Protocol: HTTPS
      SSLCertificateArns: "A"
  Resources:
    SSLCert2:
      Type: "AWS::ElasticLoadBalancingV2::ListenerCertificate"
      Properties:
        ListenerArn:
          Ref : "AWSEBV2LoadBalancerListener443"
        Certificates:
          - CertificateArn: "B"
    SSLCert3:
      Type: "AWS::ElasticLoadBalancingV2::ListenerCertificate"
      Properties:
        ListenerArn:
          Ref : "AWSEBV2LoadBalancerListener443"
        Certificates:
          - CertificateArn: "C"
  ```

* Run the build.ps1 (Windows) or build.sh (Linux) to create the 
``aws-eb-swan-demo.zip`` bundle. The bundle should contain the ``application`` 
executable compiled for Linux 64 bit, ``Procfile`` to tell Elastic Beanstalk how 
to start the application, ``.ebextensions/.config`` to configure additional
HTTPS listeners for additional domains and SSL certificates, and the ``www`` 
folder with directories for all the domains the demo will respond to. Note the
SWIFT domains used for storage do not need to be present in the www folder.

* Create a new Elastic Beanstalk Application and Environment using the 
[AWS documentation](https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/create_deploy_go.html).

* On your new Elastic Beanstalk instance set the following environment variables:
  * ``AWS_ENABLED`` - ``true``
  * ``PORT`` - ``5000``

* Give the Elastic Beanstalk Environment permissions to create and read DynamoDB
tables:
 
  * For the Elastic Beanstalk Environment role (default: `aws-elasticbeanstal-ec2-role`), attach the 
  [AmazonDynamoDBFullAccess](https://console.aws.amazon.com/iam/home#/policies/arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess$serviceLevelSummary) permission policy.

* Upload the ``aws-eb-swan-demo.zip`` bundle to the environment.

* Add an A record for each of the domains to direct traffic to the Elastic 
Beanstalk environment following the 
[AWS documentation](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-to-beanstalk-environment.html).

* The SWAN access domain will be used to sign all the outgoing Open Web IDs and
also to capture people's preferences. Register this domain with the following
URL and entering any of the details requested. This will create a record in the 
``owidcreators`` table for the domain which will contain randomly generated 
public and private signing keys.

  ```
  https://swan-access-domain/owid/register
  ```

* For each of the storage nodes that will be used for the SWIFT component of the
demo register these using the following URL. Enter the network as "swan" (no 
quotes) to match the value provided in the ``appsettings.json`` in the 
``swanNetwork`` field. Leave the others as default.

  ```
  https://swift-node-domain/swift/register
  ```

* At least one SWIFT access node is required. Repeat the previous process but 
select the "Access Node" option rather than the default "Storage Node". The 
records from these steps will be visible in the ``swiftnodes`` and 
``swiftsecrets`` tables.

* Verify the demo is working by navigating to the publisher domain. The first 
request from a web browser will result in the progress circle as SWIFT nodes
are navigated between before the preference capture page from the SWAN domain
is displayed.

## Azure App Service - with Docker

Azure App Services do not support the Go runtime natively so for deployment, Docker
containers are used. 

Azure can be used to host docker containers using Azure Container Registries. These
can then be referenced by an App Service to pull containers from.

### Pre-requisites

* Local Go version 1.15 or greater installation sufficient to run the Go command
line.

* Familiar with the concepts associated with 
[SWAN](https://github.com/SWAN-community/swan),
[SWIFT](https://github.com/SWAN-community/swift), and 
[OWID](https://github.com/SWAN-community/owid).

* Docker Engine - https://docs.docker.com/get-docker/

* Azure subscription with the ability to create App Services, Container Registries, DNS Zones and SSL Certificates

#### Azure

* Setup Azure container registry to push the demo Docker container to.
  * [Creating a container registry](https://docs.microsoft.com/en-us/azure/container-registry/container-registry-get-started-portal)
  * [Enable Admin Authentication](https://docs.microsoft.com/en-us/azure/container-registry/container-registry-authentication#admin-account)

* Setup domains to use with the demo in Azure using 
[App Service Domains](https://portal.azure.com/#blade/HubsExtension/BrowseResourceBlade/resourceType/Microsoft.DomainRegistration%2Fdomains). Domains are needed for the 
following roles.

    * SWAN access domain. For example ``51da.uk``.

    * Each of the SWIFT nodes that will support SWAN. At least two domains are 
    needed. For the purposes of the demo they may be sub domains. For example
    ``51da.uk`` and ``51db.uk``.

    * At least two different domains for publishers. For example 
    ``new-pork-limes.uk`` and ``current-bun.uk``.

    * At least two different domains for marketers. For example 
    ``cool-bikes.uk`` and ``cool-creams.uk``.

    * A domain for SSPs, Exchanges, DPSs and CMPs. For the purposes of the demo they can be sub domains. For example ``dsp.swan-demo.uk`` and ``cmp.swan-demo.uk``

* Setup SSL certificates in Azure using 
[App Service Certificates](https://portal.azure.com/#blade/HubsExtension/BrowseResource/resourceType/Microsoft.CertificateRegistration%2FcertificateOrders) 
for each of the domains.

* Create a place to store SWAN demo data using a [Storage Account](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-create?tabs=azure-portal)

#### Windows

* Set powershell unrestricted execution policy to enable `build.ps1` to execute.
  Use the following command with administrator privileges.

  `Set-ExecutionPolicy Unrestricted`

### Steps

* Get all the dependencies needed by the Go application.

  `go get -d ./...`

  `dep.ps1` can also be used in Powershell to explicitly get the dependencies 
  for this Go application.

* Add the demo domain names as folder names to the www folder. For example; the
domain `domain.com` would appear as `www/domain.com`. Alter the `config.json` 
content in each folder to indicate the purpose of the domain. This aspect of 
the demo is changing and domain examples provided with the should be reviewed 
along with the demo source code and comments to understand how multiple domains
are supported within a single demo application.

* You may need to support multiple SSL certificates if the demo deployment 
should respond to five or more domains. 
See Azure [documentation](https://docs.microsoft.com/en-us/azure/app-service/configure-ssl-certificate#import-an-app-service-certificate).

* Copy or rename ``/Dockerfile.rename`` to ``/Dockerfile`` and modify the `ENV` 
values for your environment. 

  * See [Manage storage account access keys](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-keys-manage#view-account-access-keys) on how to view and 
  retrieve access keys for your target storage account.

* Run the build.ps1 (Windows) or build.sh (Linux) to create the Docker image.
The image should contain the ``application`` executable compiled for Linux 64 
bit, the ``appsetings.json`` file and the ``www`` folder with directories for 
all the domains the demo will respond to. Note the SWIFT domains used for 
storage do not need to be present in the www folder.

* Before you can push the image to your Azure Container Registry (ACR) you will 
need to tag the image using the docker tag command. Replace <login-server> with 
the login server name of your ACR instance.

  ```
  docker tag swan-community/swan-demo-go <login-server>/swan-demo-go:latest
  ```

* Push the container to the ACR instance. 
  ```
  docker push <login-server>/swan-demo-go:latest
  ```

* Create a new App Service - see [Deploy and run a containerized web app with Azure App Service](https://docs.microsoft.com/en-us/learn/modules/deploy-run-container-app-service/5-exercise-deploy-web-app?pivots=csharp)

  * Go to the [Azure portal](https://portal.azure.com/) and select Create a resource. 

  * Search for and select Web App.

  * On the Basics tab, configure the values. The following values are for guidance:
    |Setting|Value|
    |-|-|
    |Resource Group|Select the resource group to use or create a new one.|
    |Name|e.g. swan-demo-go-app|
    |Publish|Docker Container|
    |OS|Linux|
    |Region|Select the same location that is close to you from previous exercise.|
    |App Service Plan|Use the default.|

  * Select Next: Docker >.

  * On the Docker tab, use the following values as a guide for each setting.
    |Setting|Value|
    |-|-|
    |Options|Single Container|
    |Image Source|Azure Container Registry|
    |Registry|Select your registry.|
    |Image|swan-community/swan-demo-go|
    |Tag|latest|
    |Startup Command|Leave this setting empty.|

  * Select Review and create, and then select Create.
  
* Map each of the domains to direct traffic to the App Service following the 
[Azure documentation](https://docs.microsoft.com/en-us/azure/app-service/app-service-web-tutorial-custom-domain#map-your-domain).

* The SWAN access domain will be used to sign all the outgoing Open Web IDs and
also to capture people's preferences. Register this domain with the following
URL and entering any of the details requested. This will create a record in the 
``owidcreators`` table for the domain which will contain randomly generated 
public and private signing keys. Do this for each SWAN access node and SWAN 
participant.

  ```
  https://swan-access-domain/owid/register
  ```

* For each of the storage nodes that will be used for the SWIFT component of the
demo register these using the following URL. Enter the network as "swan" (no 
quotes) to match the value provided in the ``appsettings.json`` in the 
``swanNetwork`` field. Leave the others as default.

  ```
  https://swift-node-domain/swift/register
  ```

* At least one SWIFT access node is required. Repeat the previous process but 
select the "Access Node" option rather than the default "Storage Node". The 
records from these steps will be visible in the ``swiftnodes`` and 
``swiftsecrets`` tables.

* Verify the demo is working by navigating to a publisher domain. The first 
request from a web browser will result in the progress circle as SWIFT nodes
are navigated between before the preference capture page from the SWAN domain
is displayed.

## Google Cloud Platform 

WIP - prerequisites and steps to set up the demo on GCP environment

### Pre-requisites

* Local Go version 1.15 or greater installation sufficient to run the Go command
line.

* Familiar with the concepts associated with 
[SWAN](https://github.com/SWAN-community/swan),
[SWIFT](https://github.com/SWAN-community/swift), and 
[OWID](https://github.com/SWAN-community/owid).

* Firebase project with the ability to manage Cloud Firestore.

### Google Firebase

See the Cloud FireStore [Quick Start](https://firebase.google.com/docs/firestore/quickstart)

If you are running the SWAN demo locally, then follow the steps to initialize on
your own server, otherwise, follow the steps to initialize on google cloud 
platform.

#### Steps

- Set the `GCP_PROJECT` environment variable value as your Firebase Project ID
- If running the SWAN demo locally then follow the steps to generate a [new service account](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances):
  - Create a new service account.
  - Give the new service account the `Firebase Admin SDK Administrator Service Agent` role.
  - On the service account, generate a new key.
  - Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to your key path.

## Local

### Hosts

This demo makes extensive use of multiple domains. For development purposes, 
setup the demo domains to resolve to localhost (127.0.0.1).

A sample hosts file is provided in the root of the project which contains all
domains used by the demo mapped to localhost: `hosts-sample`.

Paths to system hosts files are below.

**Windows**

In an elevated command prompt:

```cmd
notepad C:\Windows\System32\drivers\etc\hosts
```

**Linux**

As root:

```sh
nano /etc/hosts
```

### HTTPS

Some of the demo elements require HTTPS to function correctly. To enable HTTPS 
support on your local machine then you will need to add the swan demo 
certificate to your local store.

#### Windows 

`.\cert.ps1` script has been provided for convenience on windows, to run it you 
will need `openssl` available on your system path. 

On windows, git usually comes packaged with a windows `openssl` binary. This can 
be found at `%ProgramFiles%\Git\usr\bin\openssl.exe`, otherwise
`openssl` for windows can be obtained from https://wiki.openssl.org/index.php/Binaries 

#### Linux

```sh
openssl req -out uk.csr -newkey rsa:2048 -nodes -keyout uk.key -extensions req_ext -config openssl-csr.conf

openssl x509 -req -days 3650 -in uk.csr -signkey uk.key -out uk.crt -extensions req_ext -extfile openssl-csr.conf
```

### Storage

The demo features a local storage option. Creators and Swift Nodes are stored
locally in json files.

This is the default configuration of the demo. See `.vscode\launch.json.rename`
The files are stored in a folder called `.swan`, the files are named consistently
with the tables names when using cloud storage solutions.

* Creators - `.swan\creators.json`
* Swift Nodes - `.swan\swiftnodes.json`

These files are available in source control and contain pre-populated data. A 
script has been provided for convenience if you want to regenerate the data i.e. 
to use a different set of secrets.

### OWID Storage

If using multiple stores to simulate multiple stores for SWIFT, the store for 
OWID must be specified as OWID only supports one instance of store.

To specify the storage implementation to use, set the environment variable 
`OWID_STORE` to one of the following:

* Azure Table Storage - `azure`
* Google Firebase - `gcp`
* Local Storage - `local`
* AWS DynamoDB - `aws`

If the `OWID_STORE` environment is not set then the first configured store in 
the above order will be used instead.

See `.vscode/launch.json.rename` for a sample.

#### Steps:

1. delete the local storage JSON files.
2. Restart or start the SWAN demo.
3. run `.\setup-nodes.ps1`

# Notes

## Azure CosmosDB / Table Storage

If you are using Azure Storage Tables, this demo requires your storage account 
name and key to be securely stored in environment variables local to the machine 
running the demo:

If configuring an Azure App Service then see: 
[Configure an App Service app in the Azure portal](https://docs.microsoft.com/en-us/azure/app-service/configure-common#configure-app-settings)

## Visual Studio Code

Use the Command Palette (Ctrl + Shift + P) to running the 
`Go: Install/Update Tools` to install `gopkgs`, `dlv` and `gopls`.

## Prebid

The publisher new-prebid-limes shows a demo of how SWAN data might transact in Prebid.js. To run the demo you will have to build Prebid.js:

```
cd www/prebid
npm install
gulp serve-fast
```

This demo uses a demo library `new-prebid-limes.uk/swan.js` which exists simply to collect any SWAN signatures returned to the page and format them when needed for a request to the CMP info page.

The included branch of Prebid.js includes a [swanIdSystem](https://github.com/openx/Prebid.js/blob/swan-demo/modules/swanIdSystem.js) module which reads the swanId of the page from `window.swanId` and exposes it to bidders (the publisher can control which bidders receive the swanId). In the demo it is exposed as an eid:

```
{
  "user": {
    "ext": {
      "eids": [
        {
          "source": "swan.io",
          "uids": [
            {
              "id": "Am5ldy1wcmViaWQtbGltZXMudWsAVnILAEYBAAABAW5ldy1wcmViaWQtbGltZXMudWsAEAAAAO+5CvBYC06tlwUMC7LNAqACNTFkYS51awAb/wkAEAAAAKLBNyHy5UE5kcN3KQkAvpdtkGrfd/KBl3Mk1r0f7MrQKmuCpmXZqjNHFJe9pJG0a7QPHShEvYHUorxiS7RKB2w9rzP0fnC0fzeox4zE69NFAmNtcC5zd2FuLWRlbW8udWsAkkULAAIAAABvbhXF8zuDcui52ZUJqx8i5/ZvkVtqRv1H2qyq/2ixzhJ5SPEvn6RVxVJo1dX0c04ds5i6ojRAH3YbIqCUDp5a4+kCNTFkYS51awBJcgsAAAAAACJblIVvjFskyW1leVSAyTvHOQalKCCrktEPjHp9/KQHaRjbAf8/foYvRtzvI8KQhnXZqI88fGZnFyxHEOe+jCNjb29sLWNhcnMudWsNAGbf7GSQpz/Y8M/abqvDEoH3XIIxbU3qoQnVtZcbOPwcpwYZBbxyBThZsrfj76OJXe1BxtXgxdTvPb8pYFxaNbQ",
              "atype": 1
            }
          ]
        }
      ]
    }
  }
}
```

SWAN readers are required to return OWID signatures when they use SWAN data. The included [`swanBidAdapter`](https://github.com/openx/Prebid.js/blob/swan-demo/modules/swanBidAdapter.js) shows how this might be done. In the example the bidder returns a tree of OWIDs using the `BidResponse.ext.swan_owids` field in OpenRTB. The BidResponse extension is used instead of Bid extension because SWAN data can be used without the submission of a bid. These OWIDs are then passed to `swan.js` in the page.
```
BidResponse
{
  "id": "swan-bid",
  "ext": {
    "swan_owids": [
      {
        "impid": "22600961cdf2dc",
        "owid": {
          "OWID": "Am5ldy1wcmViaWQtbGltZXMudWsAVnILAEYBAAABAW5ldy1wcmViaWQtbGltZXMudWsAEAAAAO+5CvBYC06tlwUMC7LNAqACNTFkYS51awAb/wkAEAAAAKLBNyHy5UE5kcN3KQkAvpdtkGrfd/KBl3Mk1r0f7MrQKmuCpmXZqjNHFJe9pJG0a7QPHShEvYHUorxiS7RKB2w9rzP0fnC0fzeox4zE69NFAmNtcC5zd2FuLWRlbW8udWsAkkULAAIAAABvbhXF8zuDcui52ZUJqx8i5/ZvkVtqRv1H2qyq/2ixzhJ5SPEvn6RVxVJo1dX0c04ds5i6ojRAH3YbIqCUDp5a4+kCNTFkYS51awBJcgsAAAAAACJblIVvjFskyW1leVSAyTvHOQalKCCrktEPjHp9/KQHaRjbAf8/foYvRtzvI8KQhnXZqI88fGZnFyxHEOe+jCNjb29sLWNhcnMudWsNAGbf7GSQpz/Y8M/abqvDEoH3XIIxbU3qoQnVtZcbOPwcpwYZBbxyBThZsrfj76OJXe1BxtXgxdTvPb8pYFxaNbQ=",
          "Children": [
            {
              "OWID": "AnB1Ym1hdGljLnN3YW4tZGVtby51awBWcgsAAgAAAAEDLMTHhj4R3CuSKdKJb09EHlhF1zqksgAgfX3MY5PTw7u+FfC4V8PJVtBF5N1T6WWH8+LATVitImD+2COXGp6tjw==",
              "Children": [
                {
                  "OWID": "AmJpZHN3aXRjaC5zd2FuLWRlbW8udWsAVnILAAIAAAABAxX5d7Y8Pli/92GduwPACGuLrx5l+BtlHPNJKzW+7INyDyfkfbr5gMy2g5yn76HknKw5TRSZiDLU5KhO3Cx//n0=",
                  "Children": [
                    {
                      "OWID": "AmNlbnRyby5zd2FuLWRlbW8udWsAVnILAEIAAAABAGNvb2wtYmlrZXMudWsvcm9iZXJ0LWJ5ZS10RzM2cnZDZXFuZy11bnNwbGFzaC5qcGcAY29vbC1iaWtlcy51awBz6XFdFsuqNUYj4FDcMbGBUDCSrQAzw68Axx9CWnL/Hbv794nwsfqHJugwOA/y0wmuS1Uwk/9KX/NQGa+GJHj+",
                      "Children": null,
                      "Value": null
                    },
                    {
                      "OWID": "AmRhdGF4dS5zd2FuLWRlbW8udWsAVnILAEIAAAABAGNvb2wtYmlrZXMudWsvcm9iZXJ0LWJ5ZS10RzM2cnZDZXFuZy11bnNwbGFzaC5qcGcAY29vbC1iaWtlcy51awCX/oTwFpXMKFMyOQfe6mHtTuXEUkR80RQ7gjGSa4jHx/WOTRKUT4NqOuCub91+NUBvkAwpdaPy7SmS9ud2eZ0p",
                      "Children": null,
                      "Value": null
                    },
                    {
                      "OWID": "Am1lZGlhbWF0aC5zd2FuLWRlbW8udWsAVnILAEIAAAABAGNvb2wtYmlrZXMudWsvcm9iZXJ0LWJ5ZS10RzM2cnZDZXFuZy11bnNwbGFzaC5qcGcAY29vbC1iaWtlcy51awBHAac8aYTE1X+Ptg+2M5VZ2Xwcj30GYNJDfGODpsg6gnLv1cW4/QqoBhWfKOW/Im84pGi8N1E7tj/kMTtw59Dp",
                      "Children": null,
                      "Value": null
                    },
                    {
                      "OWID": "Am9hdGguc3dhbi1kZW1vLnVrAFZyCwBIAAAAAQBjb29sLWNyZWFtcy51ay9iZWUtbmF0dXJhbGxlcy11X0hqSGZrekF5TS11bnNwbGFzaC5qcGcAY29vbC1jcmVhbXMudWsAn8xrJETFajRbDn2iDMaC8/Qx/IxQRnAS8AfRO6j48dQHeCzCR4kVCc1ZC9tyEEDKbVYz0mlSsG5A/JhsJ6wnWw==",
                      "Children": null,
                      "Value": null
                    },
                    {
                      "OWID": "AnRoZXRyYWRlZGVzay5zd2FuLWRlbW8udWsAVnILAEIAAAABAGNvb2wtYmlrZXMudWsvcm9iZXJ0LWJ5ZS10RzM2cnZDZXFuZy11bnNwbGFzaC5qcGcAY29vbC1iaWtlcy51awCf58s1fIjpxQrdUHLqdNHB+3vXbxAxxXeooOZ/g2a9R0A132ZeBRcpgr2BadsoHd8JLqrt4BShkFIj0QqopMnm",
                      "Children": null,
                      "Value": null
                    },
                    {
                      "OWID": "AnpldGEuc3dhbi1kZW1vLnVrAFZyCwBIAAAAAQBjb29sLWNyZWFtcy51ay9iZWUtbmF0dXJhbGxlcy11X0hqSGZrekF5TS11bnNwbGFzaC5qcGcAY29vbC1jcmVhbXMudWsATuxHQWQ59ch+E4zA59j1eBySgeEWzPVzGd2Wn2vfYWwPcGY/9kZK82A3LWSncxrZrYO20SXD1TPiWiPf9j9A8g==",
                      "Children": [
                        {
                          "OWID": "AmxpdmVpbnRlbnQuc3dhbi1kZW1vLnVrAFZyCwACAAAAAQM55HrmAo6su1AGS8MzTn5XsPkGW32Xuq1nzaZcnXcDAK/8zNP/T5nhZyd98XUZJ13hXcQOSgcUdqRAoBXrGXLU",
                          "Children": null,
                          "Value": null
                        }
                      ],
                      "Value": -1
                    }
                  ],
                  "Value": 0
                }
              ],
              "Value": 0
            }
          ],
          "Value": null
        }
      }
    ]
  }
}
```

The swan-demo package includes a prebid route in the OpenRTB controller which is used to return SWAN BidResponses for the demo. It is an OpenRTB interface around the existing demo advertising mechanism.
