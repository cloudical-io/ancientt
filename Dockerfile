FROM golang:1.13.0-buster AS go-build

WORKDIR /go/src/app
COPY . .

RUN go get -v ./... && \
    go install -v ./...

# TODO Use fixed tag
FROM galexrt/container-toolbox:latest
LABEL maintainer="Alexander Trost <galexrt@googlemail.com>"

COPY --from=go-build /go/bin/app /bin/acntt

RUN chmod 755 /bin/acntt

ENTRYPOINT ["/bin/acntt"]

CMD ["--help"]
