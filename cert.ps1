openssl req -out uk.csr -newkey rsa:2048 -nodes -keyout uk.key -extensions req_ext -config openssl-csr.conf
openssl x509 -req -days 3650 -in uk.csr -signkey uk.key -out uk.crt -extensions req_ext -extfile openssl-csr.conf
CertUtil -addStore Root uk.crt