PROJECT_NAME=${1:-MultiProxy}

echo "linux build"
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
go build -ldflags="-w -s" -trimpath -o out/$PROJECT_NAME-$GOOS-$GOARCH

echo "windows build"
CGO_ENABLED=0
GOOS=windows
GOARCH=amd64
go build -ldflags="-w -s" -trimpath -o out/$PROJECT_NAME-$GOOS-$GOARCH

echo "darwin build"
CGO_ENABLED=0
GOOS=darwin
GOARCH=amd64
go build -ldflags="-w -s" -trimpath -o out/$PROJECT_NAME-$GOOS-$GOARCH