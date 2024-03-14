FROM --platform=${TARGETPLATFORM:-linux/amd64} debian:bookworm

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN apt-get update && apt-get install -yq --no-install-recommends \
    alien \
    build-essential \
    bzip2 \
    ca-certificates \
    cmake \
    curl \
    fonts-humor-sans \
    jq \
    git \
    gnupg \
    libaio-dev \
    libaio1 \
    libarchive-tools \
    libpq-dev \
    locales \
    locales-all \
    lsb-release \
    tzdata \
    unixodbc-dev \
    unzip \
    wget \
    zip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN echo "$TARGETARCH, $TARGETOS"
RUN if [ "$TARGETARCH" = "arm64" ]; then \
        RPM_URL=https://download.oracle.com/otn_software/linux/instantclient/instantclient-basic-linux-arm64.rpm; \
    else \
      RPM_URL=https://download.oracle.com/otn_software/linux/instantclient/2111000/oracle-instantclient-basic-21.11.0.0.0-1.x86_64.rpm; \
    fi && \
    curl -L -o /tmp/oracle-instantclient.rpm $RPM_URL && \
    alien -i /tmp/oracle-instantclient.rpm && \
    rm -rf /var/cache/yum && \
    rm -f /tmp/oracle-instantclient.rpm && \
    echo /usr/lib/oracle/21/client64/lib > /etc/ld.so.conf.d/oracle-instantclient21.conf && \
    ldconfig
ENV PATH=$PATH:/usr/lib/oracle/21/client64/bin

COPY dist /dist
RUN mkdir -p /app
RUN cp /dist/linux-*_linux_${TARGETARCH}*/knaudit-proxy /app/knaudit-proxy
RUN rm -rf /dist

CMD ["/app/knaudit-proxy", "-backend-type", "oracle"]
