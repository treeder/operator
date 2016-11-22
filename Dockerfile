FROM iron/go

WORKDIR /app
ADD VERSION /app/
COPY operator-alpine /app/operator

ENTRYPOINT ["./operator"]
