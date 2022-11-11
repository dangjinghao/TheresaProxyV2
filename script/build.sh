PROJECT_NAME=${1:-MultiProxy}
version=${2}
set CGO_ENABLED=0

oss=(windows linux darwin)
archs=(amd64 arm64)

for os in ${oss[@]}
do
for arch in ${archs[@]}
do
        if [ ${os} == "windows" ]
        then
          suffix=".exe"
        else
          suffix=""
        fi

        echo "build ${os}_${arch}"
        env GOOS=${os} GOARCH=${arch} go build -trimpath -ldflags="-s -w -X main.Version=${version}"  -o out/${PROJECT_NAME}_${os}_${arch}${suffix}

done
done