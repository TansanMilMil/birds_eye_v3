FROM mcr.microsoft.com/devcontainers/go:1.24

# Remove the broken yarn apt repository (GPG key expired)
RUN rm -f /etc/apt/sources.list.d/yarn.list \
    && apt-get update \
    && apt-get install -y vim \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
