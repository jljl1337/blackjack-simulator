# Start by building the application.
FROM golang:1.24.4 AS build

WORKDIR /go/src/app

COPY go.mod .
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o /go/bin/app cmd/blackjack-simulator/main.go

# Now copy it into our base image.
FROM scratch AS runtime

WORKDIR /app

COPY --from=build /go/bin/app /app/blackjack-simulator

ENTRYPOINT ["/app/blackjack-simulator"]