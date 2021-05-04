FROM golang:1.15-alpine as build-env

#STEP 1 build the image
#RUN apk --no-cache add git
ENV GO111MODULE=on
RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .

# Get dependenciesgo
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN GO111MODULE=on GOOS=linux CGO_ENABLED=0 go build -o sokrates cmd/main.go

FROM alpine

#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir /app
WORKDIR /app

COPY --from=build-env /app/sokrates /app/

ENV TMPDIR=/tmp
ENV GOMAXPROCS=8

EXPOSE 5000
ENTRYPOINT ["/app/sokrates"]