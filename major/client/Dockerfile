FROM golang
COPY . /go/src/github.com/joisandresky/microservices-siakadplus/major/client
WORKDIR /go/src/github.com/joisandresky/microservices-siakadplus/major/client
ENV g_major_server msp-major-service
RUN curl https://glide.sh/get | sh
RUN glide install
ENTRYPOINT go run main.go
EXPOSE 8181
