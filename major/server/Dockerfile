FROM golang
COPY . /go/src/github.com/joisandresky/microservices-siakadplus/major/server
WORKDIR /go/src/github.com/joisandresky/microservices-siakadplus/major/server
RUN curl https://glide.sh/get | sh
RUN glide install
ENTRYPOINT go run main.go
EXPOSE 9191
