FROM golang:latest
WORKDIR /app
RUN go mod init msdscourse-app
RUN go mod tidy
COPY . ./
RUN RUN CGO_ENABLED=0 GOOS=linux go build -o /msdscourse-app

ENV PORT 8080
CMD ["/msdscourse-app"]