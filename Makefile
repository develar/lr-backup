update-deps:
	go get -d -u ./...
	go mod tidy

lint:
	golangci-lint run

build-aws-functions:
	mkdir -p functions
	go get -d ./...
ifndef CLIENT_ID
	$(error CLIENT_ID is undefined)
endif
ifndef CLIENT_SECRET
	$(error CLIENT_SECRET is undefined)
endif
	rm -rf functions/*
	go build -ldflags "-s -w -X common.ClientId=${CLIENT_ID} -X common.ClientSecret=${CLIENT_SECRET} -X common.InPk=${IN_PK} -X common.InSk=${IN_SK} -X common.OutPk=${OUT_PK}" -o functions/auth cmd/auth/auth.go