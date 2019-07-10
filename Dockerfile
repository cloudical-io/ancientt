FROM fedora:30
LABEL maintainer="Alexander Trost <galexrt@googlemail.com>"

COPY packages /packages

RUN dnf -q update -y && \
    dnf --setopt=install_weak_deps=False --best install -y iperf iperf3 siege && \
    dnf clean all && \
    mkdir /workdir

USER root
WORKDIR /workdir
