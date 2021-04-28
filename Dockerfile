FROM golang
RUN go install github.com/omarelsheekh/gin-example@latest
CMD gin-example