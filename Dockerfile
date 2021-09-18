FROM ghcr.io/galexrt/container-toolbox:v20210915-101121-713

ARG BUILD_DATE="N/A"
ARG REVISION="N/A"

ARG ANCIENTT_VERSION="N/A"

LABEL org.opencontainers.image.authors="Alexander Trost <galexrt@googlemail.com>" \
    org.opencontainers.image.created="${BUILD_DATE}" \
    org.opencontainers.image.title="cloudical-io/ancientt" \
    org.opencontainers.image.description="A tool to automate network testing tools, like iperf3, in dynamic environments such as Kubernetes and more to come dynamic environments." \
    org.opencontainers.image.documentation="https://github.com/cloudical-io/ancientt/blob/main/README.md" \
    org.opencontainers.image.url="https://github.com/cloudical-io/ancientt" \
    org.opencontainers.image.source="https://github.com/cloudical-io/ancientt" \
    org.opencontainers.image.revision="${REVISION}" \
    org.opencontainers.image.vendor="cloudical-io" \
    org.opencontainers.image.version="${ANCIENTT_VERSION}"

ADD .build/linux-amd64/ancientt /bin/ancientt

RUN chmod 755 /bin/ancientt

ENTRYPOINT ["/bin/ancientt"]

CMD ["--help"]
