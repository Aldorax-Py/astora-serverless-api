build:
	set GOOS=linux
	set GOARCH=amd64
	go build -ldflags="-s -w" -o bin/main main.go

deploy_prod: build
	serverless deploy --stage prod --aws-profile astora-tech-admin