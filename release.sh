#!/bin/bash
set -ex

user="treeder"
service="operator"
version_file="main.go"
tag="latest"

./build.sh

if [ -z $(grep -Eo "[0-9]+\.[0-9]+\.[0-9]+" $version_file) ]; then
  echo "did not find semantic version in $version_file"
  exit 1
fi

perl -i -pe 's/\d+\.\d+\.\K(\d+)/$1+1/e' $version_file
version=$(grep -Eo "[0-9]+\.[0-9]+\.[0-9]+" $version_file)
echo "Version: $version"

git add -u
git commit -m "$service: $version release"
git tag -a "$version" -m "version $version"
git push
git push --tags

# Finally tag and push docker images
docker tag $user/$service:$tag $user/$service:$version

docker push $user/$service:$version
docker push $user/$service:$tag
