# Get the depencencies for this Go application.
export GOPATH=$(pwd)
go get \
github.com/Azure/azure-sdk-for-go/storage \
github.com/aws/aws-sdk-go/aws \
github.com/aws/aws-sdk-go/aws/awserr \
github.com/aws/aws-sdk-go/aws/session \
github.com/aws/aws-sdk-go/service/dynamodb \
github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute \
github.com/aws/aws-sdk-go/service/dynamodb/expression \
github.com/google/uuid \
github.com/satori/go.uuid \
github.com/bsm/openrtb \
cloud.google.com/go/firestore \
firebase.google.com/go \
google.golang.org/api/iterator \
golang.org/x/sys/unix
