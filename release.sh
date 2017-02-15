#!/bin/bash
set -ex

user="treeder"
service="operator"
version_file="main.go"
tag="latest"

# ensure working dir is clean
if [[ -z $(git status -s) ]]
then
  echo "tree is clean"
else
  echo "tree is dirty, please commit changes before running this"
  exit
fi

docker run --rm -v "$PWD":/app treeder/bump patch
version=`cat VERSION`
echo "version: $version"

./build.sh

git add -u
git commit -m "$service: $version release"
git tag -a "$version" -m "version $version"
git push
git push --tags

# Finally tag and push docker images
docker tag $user/$service:$tag $user/$service:$version

docker push $user/$service:$version
docker push $user/$service:$tag
