FROM golang:1.21

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN make build


FROM gcr.io/distroless/base-debian11

WORKDIR /
COPY --from=0 /app/bin /usr/local/bin
