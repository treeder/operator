FROM iron/go

WORKDIR /app
ADD VERSION /app/VERSION
ADD operator-alpine /app/operator

ENTRYPOINT ["./operator"]
