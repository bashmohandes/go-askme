FROM golang:1.11

WORKDIR $GOPATH/src/github.com/bashmohandes/go-askme
COPY . .
ENV GOBIN="$GOPATH/bin"
RUN go get -d ./... 
RUN go get -u github.com/gobuffalo/packr/... 
RUN go get -u github.com/lib/pq
RUN go get -u github.com/lib/pq/hstore
RUN packr install -v .
EXPOSE 8080
CMD ["go-askme"]
