FROM golang:alpine AS user
RUN apk update && apk add --no-cache git
WORKDIR /online-library-system
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/user ./cmd/user-svc
EXPOSE 8080
ENTRYPOINT ["/bin/user"]
CMD ["-logtostderr"]