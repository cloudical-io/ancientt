FROM golang:1.15.2-buster AS go-build

WORKDIR /go/src/app
COPY . .

RUN go get -v ./... && \
    go install -v ./...

FROM galexrt/container-toolbox:v20201001-123802-585
LABEL maintainer="Alexander Trost <galexrt@googlemail.com>"

COPY --from=go-build /go/bin/app /bin/ancientt

RUN chmod 755 /bin/ancientt

ENTRYPOINT ["/bin/ancientt"]

CMD ["--help"]
