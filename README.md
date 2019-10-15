# online-library-system
This is the codebase for online library management in GoLang. It contains three microservices for users, books and book-issue management. This repo is equipped with EFK logging.

The microservices can be run using the individual commands given below with the service description or they can be run collectively with the Docker.

### user-service:
This microservice accounts for user signup and login along with the authentication related tasks.

### book-service:
This microservice accounts for the books and authors related details management. It provides a structured book searching functionality with several filters and parameters.

### management-service:
This microservice is accountable for book-issue and availability management and services and maintains the record for books issue and returns.

###To run the services individually:
#### Prerequisites:
Follow the official documentation for EFK installation or click on this [link](https://docs.google.com/document/d/1s24lqsu_rhimB7s2CtGMuuweuHkmohbx2CpZ9BeaMz8/edit?usp=sharing).
- Elasticsearch
- Fluentd
- Kibana
#### user-service:
`go install ./cmd/user-svc && user-svc -logtostderr`
#### book-service:
`go install ./cmd/book-svc && book-svc -logtostderr`
#### management-service:
`go install ./cmd/management-svc && management-svc -logtostderr`

###To run the serivces using Docker:
#### Prerequisites:
- Docker
#### command:
`docker-compose up --build`

