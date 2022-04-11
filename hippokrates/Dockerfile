FROM golang:1.16-alpine

# gcc is needed for some packages
RUN apk add build-base

# add code
COPY . /app/hippokrates
WORKDIR /app/hippokrates

#needed to add the godog cli command
RUN go get github.com/cucumber/godog/cmd/godog@v0.12.0
RUN go mod download

ENTRYPOINT ["tail", "-f", "/dev/null"]
