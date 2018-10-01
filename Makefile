gen_proto:
	protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. proto/service/service.proto
	protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --python_out=plugins=grpc:. proto/service/service.proto
	protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=. proto/service/service.proto

gen_swagger:
	protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=. proto/service/service.proto
	
gen_server_tls:
	openssl req -x509 -newkey rsa:4096 -keyout server/server-key.pem -out server/server-cert.pem -days 365 -nodes -subj '/CN=localhost'

run.server:
	go run server/main.go

run.client:
	go run client/main.go

run.python:
	python proto/service/client.py