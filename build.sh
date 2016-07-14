set -ex

docker run --rm -v "$PWD":/go/src/github.com/treeder/operator -w /go/src/github.com/treeder/operator iron/go:dev go build -o operator-alpine
docker build -t treeder/operator:latest .
