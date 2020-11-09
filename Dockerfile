FROM golang:1.14-alpine as builder
RUN apk --no-cache add ca-certificates git gcc g++
WORKDIR /build/myapp

# Fetch dependencies
COPY go.mod ./
RUN go mod download

# Build
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux go build

# Create final image
FROM alpine
WORKDIR /root
COPY --from=builder /build/myapp/blog .
EXPOSE 7000
CMD ["./blog"]