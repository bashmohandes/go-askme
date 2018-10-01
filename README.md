# Go-ASKE

This GitHub repo has my attempt to build a web application in Go **without** any Web Frameworks

## Introduction
The idea is a Q/A website where users can send questions to a user they follow, the user may choose to answer all or some of these questions, the answers are public to everyone.

This is not meant to be a production ready product _at least initially_, it is meant as an educational vehicle to learn good design principals, Go language, as well as other DevOps toolchain like Docker, Kubernetes, ...etc.

## Design Goals

The code base follows Uncle Bob's [Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html) principals, where the code is broken down to

1. Entities (model folder)
2. Repositories
3. Usecases
4. Delivery

## Dependencies

This is not meant to depend on any fat frameworks, especially web frameworks, although a couple of things were used to tie things together.

1. Uber's [dig](https://go.uber.org/dig) Dependency Injection framework
2. Buffalo's [box](https://github.com/gobuffalo/packr) asset management
3. Google's [UUID](https://github.com/google/uuid) package

## Build and Run

* Install Go
* Make sure your GOPATH environment variable
* (Optional) Install Docker
* Clone the repo on your machine
* From root of the repo on your terminal run the
  following command
  ```bash
  go get -u -v ./...
 
  go run cmd/askme/main.go
  ```
* (Optional) if you prefer Docker, run the following commands
  ```bash
  docker build -t go-askme .

  docker run -it -p 8080:8080 go-askme
  ```
Then from a browser window, navigate to http://localhost:8080
