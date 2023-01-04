FROM golang:1.16-alpine
WORKDIR /server
COPY go.mod ./
COPY go.sum ./
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY ./ ./
RUN go build -o /prod
EXPOSE 8888
CMD [ "/prod" ]