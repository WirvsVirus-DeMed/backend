#!/bin/sh

createCert () {
	local c=$1
	cat > $(pwd)/${c}.sslconf <<-EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no
[req_distinguished_name]
CN = ${c}
[v3_req]
keyUsage = keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${c}
DNS.2 = DeMed-Node
EOF

	openssl genrsa -out ${c}.key 2048
	openssl req -new -sha256 -key ${c}.key -out ${c}.csr -config ${c}.sslconf
	openssl x509 -req -in ${c}.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out ${c}.crt -days 500 -sha256  -extensions v3_req -extfile ${c}.sslconf
	openssl x509 -in ${c}.crt -text -noout
	
	rm $(pwd)/${c}.sslconf
}

createCert client0
createCert client1
createCert client2
createCert client3
createCert client4
