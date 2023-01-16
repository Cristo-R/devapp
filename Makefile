test:
	KAFKA_HOST=127.0.0.1 DB_HOSTNAME=127.0.0.1 DB_USERNAME=root DB_PASSWORD=123456 DB_DATABASE=oauth2_production CGO_ENABLED=0 BASE_PATH=`pwd` go test ./... -v  -gcflags=-l -p 1  -coverprofile=coverage.out
cover:
	go tool cover -func=coverage.out

html:
	go tool cover -html=coverage.out -o coverage.html 

open:
	open coverage.html
