FROM golang:1.11

WORKDIR $GOPATH/src/github.com/bashmohandes/go-askme
COPY . .
ENV GOBIN="$GOPATH/bin"
RUN go get -d -v ./... 
RUN go get -u github.com/gobuffalo/packr/... 
RUN packr install -v .
EXPOSE 8080
CMD ["go-askme"]
