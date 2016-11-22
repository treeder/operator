Converts an RSA key to a single line with \n in between, for use in .env files. Can be used for AWS private keys. 


## Usage

```sh
docker run --rm -v $PWD:/envs -w /envs treeder/operator:rsa-key-flatten key.pem
```
