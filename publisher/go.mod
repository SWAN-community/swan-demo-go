module github.com/SWAN-community/swan-demo-go/publisher

go 1.17

require (
	github.com/SWAN-community/owid-go v0.1.2
	github.com/SWAN-community/swan-demo-go/common v0.1.3
	github.com/SWAN-community/swan-demo-go/fod v0.1.3
	github.com/SWAN-community/swan-demo-go/swanopenrtb v0.1.3
	github.com/SWAN-community/swan-go v0.1.2
	github.com/SWAN-community/swan-op-go v0.1.2
)

require (
	cloud.google.com/go v0.94.1 // indirect
	cloud.google.com/go/firestore v1.5.0 // indirect
	cloud.google.com/go/storage v1.16.1 // indirect
	firebase.google.com/go v3.13.0+incompatible // indirect
	github.com/Azure/azure-sdk-for-go v57.1.0+incompatible // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.20 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.15 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/SWAN-community/config-go v0.1.3 // indirect
	github.com/SWAN-community/salt-go v0.1.1 // indirect
	github.com/SWAN-community/swift-go v0.1.4 // indirect
	github.com/aws/aws-sdk-go v1.40.37 // indirect
	github.com/bsm/openrtb v2.1.2+incompatible // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/golang-jwt/jwt/v4 v4.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gax-go/v2 v2.1.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.8.1 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20210903162142-ad29c8ab022f // indirect
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f // indirect
	golang.org/x/sys v0.0.0-20210903071746-97244b99971b // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/api v0.56.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/SWAN-community/swan-demo-go/common => ../common

replace github.com/SWAN-community/swan-demo-go/swanopenrtb => ../swanopenrtb

replace github.com/SWAN-community/swan-demo-go/fod => ../fod

replace github.com/SWAN-community/swan-demo-go/swan-go => ../swan

replace github.com/SWAN-community/swan-demo-go/owid-go => ../owid

replace github.com/SWAN-community/swan-demo-go/swan-op-go => ../swanop
