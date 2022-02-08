FROM golang:1.16 AS builder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV TZ Asia/Shanghai

WORKDIR /app
COPY go.mod .
COPY go.sum .
#RUN GOPROXY="https://goproxy.io,direct" go mod download
RUN go mod download

COPY . .
RUN ["chmod", "+x", "/app/version.sh"]
RUN ["sh", "/app/version.sh"]

RUN go build -o agent main.go


FROM alpine AS final
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/agent /agent

CMD ["/app/agent", "server"]