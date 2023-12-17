FROM golang:1.20

WORKDIR /app

COPY go.mod .
COPY main.go .

RUN go get -u github.com/gin-gonic/gin

RUN go build -o bin .

EXPOSE 8000

CMD ["./bin"]