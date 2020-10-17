FROM golang:1.14
WORKDIR /go/src/app
COPY . /go/src/app
RUN go get -d -v github.com/gorilla/mux
RUN go get -d -v github.com/google/uuid
RUN go get -d -v github.com/jinzhu/gorm
RUN go get -d -v github.com/mattn/go-sqlite3
CMD ["go","run","main.go"]