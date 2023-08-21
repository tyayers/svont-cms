############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /cmsservice .

############################
# STEP 2 build a small image
############################
FROM gcr.io/distroless/static-debian11
#FROM gcr.io/distroless/cc

ENV GIN_MODE=release

COPY --from=builder /cmsservice .

ENTRYPOINT ["/cmsservice"]

EXPOSE 8080
