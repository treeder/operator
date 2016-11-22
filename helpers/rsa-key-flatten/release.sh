set -ex

./build.sh

docker push treeder/operator:rsa-key-flatten
