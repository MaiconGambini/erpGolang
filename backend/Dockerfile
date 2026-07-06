FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/goerp-api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/goerp-api /goerp-api
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/goerp-api"]
