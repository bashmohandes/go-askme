# Go-ASKME [![Build Status](https://travis-ci.org/bashmohandes/go-askme.svg?branch=master)](https://travis-ci.org/bashmohandes/go-askme)

This GitHub repo has my attempt to build a web application in Go **without** any Web Frameworks

## Introduction
The idea is a Q/A website where users can send questions to a user they follow, the user may choose to answer all or some of these questions, the answers are public to everyone.

This is not meant to be a production ready product _at least initially_, it is meant as an educational vehicle to learn good design principals, Go language, as well as other DevOps toolchain like Docker, Kubernetes, ...etc.

## Design Goals

The code base follows Uncle Bob's [Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html) principals, where the code is broken down to

1. Entities (model folder)
2. Components
    1. Questions
    2. Answers
    3. Users
3. Use Cases
4. Shared

### Entities
The main domain models, which in our case _so far_ **Question**, **Answer**, **User** as well as base entity types and helpers, and basic domain operations on domain models using the domain ubiquitous language, like Ask, Answer, Like ... etc

### Components
Each component consists of all the basic layers needed to complete this component from top to bottom, like **Use Cases**, **Repositories**, **Tests**

### Framework
This is shaping up to be a tiny MVC framework, interesting

## Dependencies

This is not meant to depend on any fat frameworks, especially web frameworks, although a couple of things were used to tie things together.

1. Uber's [dig](https://go.uber.org/dig) Dependency Injection framework
2. Buffalo's [box](https://github.com/gobuffalo/packr) asset management
3. Google's [UUID](https://github.com/google/uuid) package
4. joho/godotenv [godotenv](https://github.com/joho/godotenv)

## Build and Run

* Install Go
* Make sure your GOPATH environment variable
* (Optional) Install Docker
* Clone the repo on your machine
* From root of the repo on your terminal run the
  following command
  ```bash
  go get -u -v ./...

  docker run --env-file=.env -p 5432:5432 --rm postgres:latest
 
  go run main.go
  ```
* (Optional) if you prefer Docker, run the following commands
  ```bash
  docker build -t go-askme .

  docker run --env-file=.env --rm -p 8080:8080 go-askme
  ```
Then from a browser window, navigate to http://localhost:8080
