FROM golang:1.20-alpine as build

WORKDIR /go/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only 
# redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app /
EXPOSE 3000
CMD ["/app"]
