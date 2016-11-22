$ErrorActionPreference = "Stop"

$user = "treeder"
$image = "operator"
$tag = "env-to-args"

# Bump version
$version = Get-Content VERSION
Write-Host "before: " $version
docker run --rm -v ${pwd}:/app treeder/bump patch
$version = Get-Content VERSION
Write-Host "after: " $version

docker build -t $user/${image}:${tag} .

git add -u
git commit -m "${tag}: $version release"
git tag -a "$tag-$version" -m "$tag version $version"
git push
git push --tags

# Finally tag and push docker images
docker tag $user/${image}:$tag $user/${image}:${tag}-$version

docker push $user/${image}:${tag}-$version
docker push $user/${image}:$tag
