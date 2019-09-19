FROM golang
COPY . /go/src/github.com/joisandresky/microservices-siakadplus/course/server
WORKDIR /go/src/github.com/joisandresky/microservices-siakadplus/course/server
RUN curl https://glide.sh/get | sh
RUN glide install
ENTRYPOINT go run main.go
EXPOSE 9292
