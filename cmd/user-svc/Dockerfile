#FROM golang:alpine AS user
#RUN apk update && apk add --no-cache git
#WORKDIR /library-management
#COPY ./go.mod ./go.sum ./
#RUN go mod download
#COPY . .
#RUN go build -o ./bin/user ./cmd/user-svc
#EXPOSE 8000
#COPY ./bin .
#ENTRYPOINT ["./bin/user-server"]
#

FROM alpine
WORKDIR /library-management
COPY ./bin .
ENTRYPOINT ["./user-server"]

