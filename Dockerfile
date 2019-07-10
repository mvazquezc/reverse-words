FROM golang:latest
WORKDIR /go/src/github.com/mvazquezc/reverse-words/
COPY main.go .
RUN go get github.com/gorilla/mux && go get github.com/prometheus/client_golang/prometheus/promhttp
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch
COPY --from=0 /go/src/github.com/mvazquezc/reverse-words/main .
EXPOSE 8080
CMD ["/main"]
