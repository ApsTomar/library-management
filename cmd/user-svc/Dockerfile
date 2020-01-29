FROM golang:alpine AS user
RUN apk update && apk add --no-cache git
WORKDIR /library-management
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/user ./cmd/user-svc
EXPOSE 8000
ENTRYPOINT ["/bin/user"]
