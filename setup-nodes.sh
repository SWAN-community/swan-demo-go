#!/bin/bash
NETWORK="swan"
NODES=( \
    "51da.uk" \
    "51db.uk" \
    "51dc.uk" \
    "51dd.uk" \
    "51de.uk" \
)
EXPIRYDATE=$(date -d '+90 days' '+%Y-%m-%d')

echo "Network: ${NETWORK}"
echo "Expiry: ${EXPIRYDATE}"

read -r -p "Ok? [y/N] " response
if ! [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    exit 0
fi

## Set-up SWIFT access nodes as OWID creators
for n in ${NODES[*]}; do
    echo "51Degrees : ${n}" 
    URL="http://${n}/owid/register?name=51Degrees"
    echo $URL
    RES=$(curl --write-out '%{http_code}' --silent --output /dev/null ${URL})
    echo $RES
done

## Set-up SWAN participant as OWID creators
DIR=$(find www -maxdepth 1 -mindepth 1 -type d)
for d in $DIR; do
    if [ -e "${d}/config.json" ]; then
        NAME=$(jq --raw-output "if .Name != null then .Name else .name end" "${d}/config.json" | tr " " +)
        DOMAIN=$(basename ${d})
        echo "${NAME} : ${DOMAIN}" 
        URL="http://${DOMAIN}/owid/register?name=${NAME}"
        echo $URL
        RES=$(curl --write-out '%{http_code}' --silent --output /dev/null ${URL})
        echo $RES
    fi
done

## Set-up SWIFT access nodes
for n in ${NODES[*]}; do
    echo "51Degrees : ${n}" 
    URL="http://${n}/swift/register?network=${NETWORK}&expires=${EXPIRYDATE}&role=0"
    echo $URL
    RES=$(curl --write-out '%{http_code}' --silent --output /dev/null ${URL})
    echo $RES
done

# ## Set-up SWIFT storage Nodes
for n in ${NODES[*]}; do
    for i in `seq 1 30`; do
        echo "51Degrees : ${i}.${n}" 
        URL="http://${i}.${n}/swift/register?network=${NETWORK}&expires=${EXPIRYDATE}&role=1"
        echo $URL
        RES=$(curl --write-out '%{http_code}' --silent --output /dev/null ${URL})
        echo $RES
    done
done
