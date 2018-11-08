FROM golang:1.11

ENV GOBIN="$GOPATH/bin"
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code from the context.
COPY . .
RUN go get -u github.com/gobuffalo/packr/...
RUN packr install -v .
EXPOSE 8080
CMD ["go-askme"]
