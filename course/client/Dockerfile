FROM golang
COPY . /go/src/github.com/joisandresky/microservices-siakadplus/course/client
WORKDIR /go/src/github.com/joisandresky/microservices-siakadplus/course/client
ENV g_course_server msp-course-service
RUN curl https://glide.sh/get | sh
RUN glide install
ENTRYPOINT go run main.go
EXPOSE 8181
