# Operator

![Operator](https://tctechcrunch2011.files.wordpress.com/2015/04/matrix-operator.jpg)

## Configuration

Create a file called `.env`, copy and paste the following into it and fill in the blanks. 

```sh
# For private images on DockerHub:
DOCKER_USERNAME=A
DOCKER_PASSWORD=B

# For deployment to EC2
AWS_ACCESS_KEY=X
AWS_SECRET_KEY=Y
AWS_KEY_PAIR=X # name of keypair
AWS_PRIVATE_KEY= # For your private key, wrap it with double quotes and replace all the new lines with `\n` values so they can go in one line, including one at the very end, eg: "...uNQgmDXEbU\n-----END RSA PRIVATE KEY-----\n". Docker's --env-file does not work with multi-line values. 
AWS_SUBNET_ID=X
AWS_SECURITY_GROUP=X
        
# For streaming logs to a syslog service, if set logspout will be installed too:
SYSLOG_URL=logs.papertrail.com:1234
```

## Commands:

### Deploy

Deploy your image to a new server, or if it's already on a server, it will just update

```
docker run --rm -it --env-file .env treeder/operator --name myapp -e X=Y IMAGE
```

### List instances

```
docker run --rm -it --env-file .env treeder/operator --name myapp instances
```

### Scale out

# TODO: 

Add servers to an app cluster

```sh
docker run --rm -it --env-file .env treeder/operator deploy --add --name someapp -e X=Y IMAGE
```

### Run SSH command on all instances of an app

Run an SSH command across all instances of an app

```sh
docker run --rm -it --env-file .env treeder/operator sh --name myapp 'some ssh command'
# examples:
docker run --rm -it --env-file .env treeder/operator sh --name myapp 'docker ps'
docker run --rm -it --env-file .env treeder/operator sh --name functions 'docker logs myapp'
```

The server will be tagged with app name for redeploying. 

TODO: add servers to a load balancer, choose instance size, maybe a yaml file?
TODO: how to set instance type
TODO: Add logspout if SYSLOG_URL specified

```yml
image: treeder/hello-sinatra
env_vars: Can put env vars in here?  But probably shouldn't have them in source control anyways
load_balancer: NAME
instance_type: x.large
dns: some cloudflare info here?  keys in env vars
monitoring: logspout, etc
```

Support Docker Compose too. 

## Utility containers

### Base64 Google Cloud credentials for use in .env files 

```sh
docker run --rm -v $PWD:/envs -w /envs treeder/operator:google-creds-flatten google-creds.json > creds.tmp
```

Then take the output in creds.tmp and put into your .env file. 

### Convert .env files into -e params when you don't want to deal with moving around .env files

eg: 

```sh
docker run --rm $(docker run --rm -v $PWD:/envs -w /envs treeder/operator:env-to-args .env) -p 8080:8080 IMAGE
```
