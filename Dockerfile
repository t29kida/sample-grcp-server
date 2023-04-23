FROM golang:1.20-bullseye as dev

ENV ROOT=/app

RUN apt update \
    && apt clean \
    && rm -rf /var/lib/opt/lists/*

WORKDIR $ROOT

COPY go.mod go.sum ./

RUN go mod download

EXPOSE 8080

FROM golang:1.20-bullseye as builder

ENV ROOT=/app

RUN apt update \
    && apt clean \
    && rm -rf /var/lib/opt/lists/*

WORKDIR $ROOT

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./backend cmd/main.go

RUN chmod +x /app/backend

FROM scratch as prod

ENV ROOT=/app

WORKDIR ${ROOT}

COPY --from=builder ${ROOT}/backend ${ROOT}

CMD ["/app/backend"]

EXPOSE 8080
