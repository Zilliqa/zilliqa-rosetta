FROM golang:1.14 as rosetta-build-stage
WORKDIR /app
COPY . ./
RUN go build main.go
RUN ls


FROM ubuntu:18.04
WORKDIR /rosetta
COPY --from=rosetta-build-stage /app/main /rosetta/main
COPY --from=rosetta-build-stage /app/config.local.yaml /rosetta/config.local.yaml
EXPOSE 8080
ENTRYPOINT ["/rosetta/main"]
