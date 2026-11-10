FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/service-plan .

FROM alpine:3.20
RUN adduser -D -u 10001 app
USER app
COPY --from=build /out/service-plan /usr/local/bin/service-plan
ENV GIN_MODE=release
EXPOSE 8081
ENTRYPOINT ["service-plan"]
