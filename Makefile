.PHONY: help

help:
	@echo "Usage:"
	@echo "  make build              compile the Go code."
	@echo "  make build-stealth      compile the Go code in stealth mode, resulting no windows."
	@echo "  make push               push the zip file to the production server."
	@echo "  make zip                zip the executable and password-protect it."

build:
	GOOS=windows GOARCH=amd64 go build -o shell.exe shellMenu.go

build-stealth:	
	GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o shell.exe shellMenu.go

zip:
	zip -P coucou shell.zip shell.exe

push:
	scp shell.zip prod.redteams:~
