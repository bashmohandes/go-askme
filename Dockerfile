FROM golang:1.11

ENV GOBIN="$GOPATH/bin"
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code from the context.
COPY . .

# clean any windows file endings
RUN apt-get update && apt-get install -y dos2unix
RUN find . -type f -exec dos2unix {} \;

# install with packr to embed resources
RUN go get github.com/gobuffalo/packr/...
RUN packr install -v .

EXPOSE 8080
CMD [ "go-askme" ]
