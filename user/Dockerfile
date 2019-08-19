FROM golang
WORKDIR /go/src/github.com/joisandresky/microservices-siakadplus/user

ADD . ./

RUN curl https://glide.sh/get | sh
RUN glide install

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .

EXPOSE 8181

CMD ["go", "run", "main.go"]
