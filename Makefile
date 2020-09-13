APP_VARS = -X github.com/develar/lr-backup/pkg/common.clientId=${CLIENT_ID}\
-X github.com/develar/lr-backup/pkg/common.clientSecret=${CLIENT_SECRET}\
-X github.com/develar/lr-backup/pkg/common.InPk=${IN_PK}\
-X github.com/develar/lr-backup/pkg/common.InSk=${IN_SK}\
-X github.com/develar/lr-backup/pkg/common.OutPk=${OUT_PK}

APP_NAME = Lightroom Backup

update-deps:
	go get -d -u ./...
	go mod tidy

lint:
	golangci-lint run

build-aws-functions: check-env
	mkdir -p functions
	go get -d ./...

	rm -rf functions/*
	go build -ldflags "-s -w ${APP_VARS}" -o functions/auth cmd/auth/auth.go

build: check-env
	go get -d ./...

	go build -ldflags "-s -w ${APP_VARS}" -o dist/lr-backup ./cmd/lr-backup


# go get fyne.io/fyne/cmd/fyne
# https://iconscout.com/icon-editor is a good site to draw icon
package: build
	fyne package -os darwin -executable ./dist/lr-backup -icon ./resources/lr-backup.png -appID org.develar.lr-backup -name '${APP_NAME}'
	mv "${APP_NAME}.app" "dist/${APP_NAME}.app"
	cp resources/Info.plist "dist/${APP_NAME}.app/Contents/Info.plist"

check-env:
ifndef CLIENT_ID
	$(error CLIENT_ID is undefined)
endif
ifndef CLIENT_SECRET
	$(error CLIENT_SECRET is undefined)
endif