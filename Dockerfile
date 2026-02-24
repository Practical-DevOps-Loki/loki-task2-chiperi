FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/webapp ./main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /out/webapp /app/webapp
COPY --from=build /src/public /app/public
ENV PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/webapp"]
