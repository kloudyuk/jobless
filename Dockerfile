FROM golang:1.15 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o jobless

FROM alpine:3.13
COPY --from=build /build/jobless /usr/bin/jobless
CMD ["jobless"]
