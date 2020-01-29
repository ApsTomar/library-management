FROM golang:alpine AS management
RUN apk update && apk add --no-cache git
WORKDIR /library-management
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/management ./cmd/management-svc
EXPOSE 8002
ENTRYPOINT ["/bin/management"]
