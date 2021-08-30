FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hashService /app/cmd/server/server.go

FROM scratch
COPY --from=builder /app/configs /usr/bin/configs
COPY --from=builder /app/hashService /usr/bin/hashService
ENTRYPOINT [ "/usr/bin/hashService" ]