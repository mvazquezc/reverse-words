FROM docker.io/library/golang:1.22
WORKDIR /go/src/github.com/mvazquezc/reverse-words/
COPY main.go .
COPY go.mod .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch
COPY --from=0 /go/src/github.com/mvazquezc/reverse-words/main .
EXPOSE 8080
USER 9999
CMD ["/main"]
