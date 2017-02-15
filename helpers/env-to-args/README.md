Converts .env files to Docker env var params. 

## Usage

### Docker

For Docker, you can use it like this:

```sh
docker run --rm $(docker run --rm -v $PWD:/envs -w /envs treeder/operator:env-to-args .env) -p 8080:8080 treeder/hello
```

### For Bash

This one exports your env vars and is used like so:



```sh
eval $( --rm -v $PWD:/envs -w /envs treeder/operator:env-to-args .env --type sh)
# then start your program
./hello
```

### Powershell

For Powershell, you want to set you env vars, then run your program, so use it in a ps1 script like this:

```ps1
$ec1 = "docker run --rm -v ${HOME}/configs/bots.haus:/envs -w /envs treeder/operator:env-to-args dev.env --type ps"
$envs = iex $ec1
iex "$envs"
# start your program here
./hello
```

## Building this image

Run release.sh or release.ps1.
