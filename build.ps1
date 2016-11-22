$ErrorActionPreference = "Stop"

$username="treeder"
$image="operator"

$cmd = $args[0]
Write-Host "cmd: $cmd"

function quick($args2) {
    go build
    # could allow user to pass in .env file here
    # If ($args2.Count -ge 1) {
    #     $env_file = "$($args2[1])"
    # }
    # $env:TestVariable = "This is a test environment variable."
    $ec1 = "docker run --rm -v ${pwd}:/envs -w /envs treeder/operator:env-to-args .env --type ps"
    $envs = iex $ec1
    iex "$envs"
    Get-ChildItem Env:
    ./operator
}

function build () {
    docker run --rm -v "${pwd}":/go/src/github.com/treeder/operator -w /go/src/github.com/treeder/operator iron/go:dev go build -o operator-alpine
    docker build -t ${username}/${image}:latest .
}

function release() {
    # ensure tree is clean
    # http://stackoverflow.com/questions/13738634/how-can-i-check-if-a-string-is-null-or-empty-in-powershell
    if(git status -s) {
        echo "tree is dirty, commit changes before releasing."
        exit
    }

    # bump version
    docker run --rm -v ${pwd}:/app treeder/bump patch
    $version = Get-Content .\VERSION -Raw 
    echo "version: $version"

    build 

    # tag it
    git add -u
    git commit -m "version $version"
    git tag -a "$version" -m "version $version"
    git push
    git push --tags

    docker tag $username/${image}:latest $username/${image}:${version}

    # push it
    docker push $username/${image}:latest
    docker push $username/${image}:${version}
}


switch ($cmd)
{
    "quick" { quick($args) }
    "build" { build }
    "release" {release}
    default {"Invalid command: $cmd"}
}
