FROM iron/go

WORKDIR /app
COPY operator-alpine /app/operator

ENTRYPOINT ["./operator"]
