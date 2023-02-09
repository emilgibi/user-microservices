FROM golang:1.19

WORKDIR /app 

COPY . . 

RUN go build main.go

EXPOSE 8083

CMD ["./main"]