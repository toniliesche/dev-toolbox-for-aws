# MIT License
# Copyright (c) 2025 Toni Liesche
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

FROM amazonlinux:2 AS builder

ARG TARGETARCH
RUN yum install -y gcc gcc-c++ git make tar gzip && \
  curl -OL https://go.dev/dl/go1.24.0.linux-${TARGETARCH}.tar.gz && \
  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.0.linux-${TARGETARCH}.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
ARG CACHEBUST
RUN echo ${CACHEBUST} && \
    git clone https://github.com/toniliesche/dev-toolbox-for-aws.git /tmp/dev-toolbox-for-aws

WORKDIR /tmp/dev-toolbox-for-aws

ARG BUILDVER
RUN --mount=type=cache,target=/root/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    git checkout ${BUILDVER} && \
    make build-secrets-manager

FROM amazonlinux:2

ARG BUILDVER
LABEL org.opencontainers.image.authors="Toni Liesche <docker@toniliesche.de>"
LABEL org.opencontainers.image.url="https://hub.docker.com/r/tliesche/secrets-manager"
LABEL org.opencontainers.image.documentation="https://raw.githubusercontent.com/toniliesche/dev-toolbox-for-aws/refs/heads/develop/README.md"
LABEL org.opencontainers.image.source="https://github.com/toniliesche/dev-toolbox-for-aws"
LABEL org.opencontainers.image.version=${BUILDVER}
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="Local Secrets Manager Mock"
LABEL org.opencontainers.image.description="Custom AWS compatible Secrets Manager for Local Development"

COPY --from=builder /tmp/dev-toolbox-for-aws/build/secrets-manager /usr/local/bin/secrets-manager

CMD [ "/usr/local/bin/secrets-manager" ]
