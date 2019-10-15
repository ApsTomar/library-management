FROM golang:alpine AS book
RUN apk update && apk add --no-cache git
WORKDIR /online-library-system
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/book ./cmd/book-svc
EXPOSE 8081
ENTRYPOINT ["/bin/book"]
CMD ["-logtostderr"]