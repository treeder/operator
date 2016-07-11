

Env vars required:

```sh
AWS_ACCESS_KEY=X
AWS_SECRET_KEY=Y
AWS_PEM=`cat yourpemfile.pem`
# for private images on DockerHub:
DOCKER_USERNAME=A
DOCKER_PASSWORD=B
```

Commands:

```
# Deploy your image to a new server, or if it's already on a server, it will just update
dj deploy -e X=Y IMAGE
```

TODO: add servers to a load balancer, choose instance size, maybe a yaml file?

```yml
image: treeder/hello-sinatra
load_balancer: NAME
instance_type: x.large
dns: some cloudflare info here?  keys in env vars
```
