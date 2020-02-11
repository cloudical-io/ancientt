FROM golang:1.13.0-buster AS go-build

WORKDIR /go/src/app
COPY . .

RUN go get -v ./... && \
    go install -v ./...

FROM galexrt/container-toolbox:v20200211
LABEL maintainer="Alexander Trost <galexrt@googlemail.com> and Michal Janus <michal.janus@cloudical.io>"

COPY --from=go-build /go/bin/app /bin/ancientt

RUN chmod 755 /bin/ancientt

ENTRYPOINT ["/bin/ancientt"]

CMD ["--help"]
