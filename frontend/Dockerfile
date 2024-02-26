FROM golang:latest
WORKDIR /app
RUN go mod init msdscourse-app
RUN go mod tidy
COPY . ./
RUN go build -o /msdscourse-app

ENV PORT 8080
CMD ["/msdscourse-app"]