set -ex

./build.sh

docker push treeder/operator:google-creds-flatten
