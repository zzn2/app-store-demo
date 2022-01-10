FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN mkdir /src
COPY src /src
WORKDIR /src
RUN go build -o main .


FROM scratch
COPY --from=builder /src/main /
ENTRYPOINT ["/main"]
