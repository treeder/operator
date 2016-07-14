

Env vars required:

```sh
# For deployment to EC2
AWS_ACCESS_KEY=X
AWS_SECRET_KEY=Y
AWS_KEY_PAIR=X # name of keypair
AWS_PRIVATE_KEY=`cat yourpemfile.pem`
AWS_SUBNET_ID=X
AWS_SECURITY_GROUP=X
        
# For private images on DockerHub:
DOCKER_USERNAME=A
DOCKER_PASSWORD=B
# For streaming logs to a syslog service (auto installs logspout):
SYSLOG_URL=udp://papertrail.com:1234
```

Commands:

```
# Deploy your image to a new server, or if it's already on a server, it will just update
./operator deploy --name someapp -e X=Y IMAGE

# Add servers to an app cluster
./operator deploy --add --name someapp -e X=Y IMAGE

# Run an SSH command across all instances of an app
./operator sh --name someapp 'some ssh command' 
```

The server will be tagged with app name for redeploying. 

TODO: add servers to a load balancer, choose instance size, maybe a yaml file?
TODO: how to set instance type

```yml
image: treeder/hello-sinatra
env_vars: Can put env vars in here?  But probably shouldn't have them in source control anyways
load_balancer: NAME
instance_type: x.large
dns: some cloudflare info here?  keys in env vars
monitoring: logspout, etc
```

Support Docker Compose too. 
