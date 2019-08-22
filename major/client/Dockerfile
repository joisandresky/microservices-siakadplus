FROM golang
COPY . /go/src/github.com/joisandresky/microservices-siakadplus/major/client
WORKDIR /go/src/github.com/joisandresky/microservices-siakadplus/major/client
RUN curl https://glide.sh/get | sh
RUN glide install
ENTRYPOINT go run main.go
EXPOSE 8181
