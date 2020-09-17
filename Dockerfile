FROM golang:alpine as builder

WORKDIR /

COPY . .
RUN CGO_ENABLED=0 go build sessionservice.go

FROM scratch
COPY --from=builder /sessionservice /
ENTRYPOINT ["/sessionservice"]
