FROM golang:1.24 AS build

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o gry -a -installsuffix cgo

FROM alpine:3

RUN addgroup -S gry && adduser -S gry -G gry

WORKDIR /opt/app

COPY --from=build /opt/app/gry .

RUN chown gry:gry gry
USER gry

CMD ["./gry"]
