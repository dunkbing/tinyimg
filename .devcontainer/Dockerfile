FROM --platform=linux/amd64 mcr.microsoft.com/vscode/devcontainers/base:0-bullseye

ENV DENO_INSTALL=/deno
RUN mkdir -p /deno \
    && curl -fsSL https://deno.land/x/install/install.sh | sh \
    && chown -R vscode /deno

ENV PATH=${DENO_INSTALL}/bin:${PATH} \
    DENO_DIR=${DENO_INSTALL}/.cache/deno

# install go
RUN curl -OL https://golang.org/dl/go1.21.5.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz \
    && rm go1.21.5.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# install flyctl
RUN curl -L https://fly.io/install.sh | sh

RUN apt update && apt upgrade -y
RUN apt install pngquant jpegoptim imagemagick webp python3 python3-pip -y
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin
RUN pip install oci-cli --upgrade
