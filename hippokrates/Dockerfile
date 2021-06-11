FROM golang:1.15

# gcc is needed for some packages
#RUN apk add build-base
COPY . /app/hippokrates
WORKDIR /app/hippokrates

#needed to add the godog cli command
RUN go get github.com/cucumber/godog/cmd/godog@v0.11.0
RUN go mod download

ENTRYPOINT ["tail", "-f", "/dev/null"]
