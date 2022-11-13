FROM golang:alpine3.16 as builder
RUN apk update && apk upgrade && apk add build-base git make sed
RUN go install github.com/silenceper/gowatch@latest

#swagger addition via swagger-ui inyection
WORKDIR /go/src/github.com/eduardohoraciosanto/users/swagger
COPY ./oas/oas.yml ./swagger.yml
RUN git clone https://github.com/swagger-api/swagger-ui && \
  cp -r swagger-ui/dist/. . && rm -r swagger-ui/ && sed -i 's+https://petstore.swagger.io/v2/swagger.json+/swagger/swagger.yml+g' index.html


WORKDIR /go/src/github.com/eduardohoraciosanto/users
COPY . .

RUN go mod tidy
RUN GIT_COMMIT=$(git rev-parse --short HEAD) && \
  go build -o service -ldflags "-X 'github.com/eduardohoraciosanto/users/internal/config.serviceVersion=$GIT_COMMIT'"

FROM alpine:3.13 

COPY --from=builder /go/src/github.com/eduardohoraciosanto/users/service /
COPY --from=builder /go/src/github.com/eduardohoraciosanto/users/swagger /swagger

ENTRYPOINT [ "./service" ]