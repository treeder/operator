Converts .env files to Docker env var params. 

## Usage

```sh
docker run --rm -v $PWD:/envs -w /envs treeder/operator:env-to-args .env
```
