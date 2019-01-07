FROM golang:latest
ENV PORT  8080

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/gorilla/mux

RUN go build -o main.go .
EXPOSE 8080
CMD /app/main.go
