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
5. [gorm](github.com/jinzhu/gorm)
6. [Okta Jwt Verifier](github.com/okta/okta-jwt-verifier-golang)

## Build and Run

* Install Go
* Make sure your GOPATH environment variable
* (Optional) Install Docker
* Clone the repo on your machine
* Create a new .env file using the .env.dist as a template
    * Fill in the missing secrets suitable for your environment
    * Never check in the .env file, it is already included in .gitignore, and never add actual secrets in .env.dist
* There are two authentications methods supported
    * Basic Username / Password Auth
        * if you want to use this mode, make sure that AuthController is used in the ask.go instead of OktaController
    * [Okta](https://developer.okta.com) based auth
        * You will need to sign up for Okta developer account, more information on setting up mentioned below
* From root of the repo on your terminal run the
  following command
  ```bash
  docker-compose -f docker-local.yml up
 
  go run main.go
  ```
* (Optional) if you prefer Docker, run the following commands
  ```bash
  docker build -t go-askme .

  docker run --env-file=.env --rm -p 8080:8080 go-askme
  ```
* (Optional) Use docker compose
  ```bash
  docker-compose build

  docker-compose up
  ```
  
Then from a browser window, navigate to http://localhost:8080


# Setting up Okta Developer Account
For more up to date information make sure to read through [Okta developer docs](https://developer.okta.com/use_cases/authentication/)

1. Sign up for a new Okta developer account
2. Create a new Application
    1. Choose Web as an application type
3. Add a login redirect URI
    1. For local development add
    http://localhost:8080/authorization-code/callback
4. Add a logout redirect URI
    1. For local development add
    http://localhost:8080	
5. Accept other defaults, then Save
6. Copy generated ClientID, Client Secret, and set them to corresponding OKTA_CLIENT_ID, OKTA_CLIENT_SECRET in your .env file
7. Navigate to API > Authorization Servers
    1. Copy the default Authorization Server, and set it to OKTA_ISSUER in your .env file
    2. Click on the *Trusted Origins* tab
        1. Add a new origin, for local development add http://localhost:8080 as both a Redirect & CORS
8. To enable Registration, navigate to Users > Registration
    1. Click on Enable Registration, fill in the details you need
    2. Make sure to set the *Default redirect* to http://localhost:8080

