FROM iron/go

WORKDIR /app
ADD VERSION .
COPY operator-alpine /app/operator

ENTRYPOINT ["./operator"]
