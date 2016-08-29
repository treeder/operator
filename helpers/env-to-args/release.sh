set -ex

./build.sh

docker push treeder/operator:env-to-args
