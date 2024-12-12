# Use an official minimal Linux distribution as the base image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Install necessary packages and dependencies
# Replace 'your-packages' with actual package names
RUN apk update && \
    apk add --no-cache \
    ca-certificates \
    openssl \
    bash \
    git \
    gcc \
    musl-dev \
    linux-headers \
    make \
    # Blockchain-specific dependencies
    leveldb \
    libsodium \
    libsecp256k1-dev \
    # Add any other required packages here \
    && rm -rf /var/cache/apk/*

# Set environment variables if necessary
ENV SOME_ENV_VAR=some_value

# Copy any additional files or scripts if needed
# COPY ./your-additional-files /app/

# Install Go (if needed for building packages)
ENV GOLANG_VERSION 1.18.6

RUN wget -q https://dl.google.com/go/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz && \
    rm go$GOLANG_VERSION.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# Copy your Go modules files to cache dependencies
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy your application source code
COPY . .

# Build your packages or libraries if necessary
# For example, build a shared library
# RUN make && make install

# Optional: Clean up to reduce image size
RUN rm -rf /usr/local/go/pkg/* && \
    rm -rf /root/.cache

# Set the entrypoint or command if this image is meant to be run
# CMD ["your-command"]
