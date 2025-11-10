# Stage 1: Builder
FROM golang:1.24.9-bookworm AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy specific directories/files needed for the build
# We copy go.mod and go.sum first to leverage Docker's caching if they haven't changed
COPY go.mod .
COPY go.sum .
# Download dependencies
RUN go mod download
RUN go mod tidy

# Copy the specific source code directories, exclude dll and .exe file
COPY . .
RUN rm embeddings.exe

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o embeddings .


# Stage 2: Production
# google example showed this image as base, i dont know enough to know if this is good or not for this use case
FROM debian:bookworm-slim

# Set the working directory
WORKDIR /app

# Install curl and tar for fetching ONNX
RUN apt-get update && apt-get install -y curl tar && rm -rf /var/lib/apt/lists/*

# Fetch ONNX
# https://github.com/microsoft/onnxruntime/releases/download/v1.22.0/onnxruntime-linux-x64-1.22.0.tgz
ARG GITHUB_USER="microsoft"
ARG GITHUB_REPO="onnxruntime"
ARG RELEASE_TAG="v1.22.0" # Or specify a specific release tag
ARG TGZ_FILENAME="onnxruntime-linux-x64-1.22.0.tgz" # The name of your .tgz file in the release assets

RUN curl -L "https://github.com/${GITHUB_USER}/${GITHUB_REPO}/releases/download/${RELEASE_TAG}/${TGZ_FILENAME}" -o /tmp/${TGZ_FILENAME} \
    && mkdir -p /tmp/extracted \
    && tar -xzf /tmp/${TGZ_FILENAME} -C /tmp/extracted \
    && cp /tmp/extracted/onnxruntime-linux-x64-1.22.0/lib/libonnxruntime.so* /usr/local/lib/ \
    && cp /tmp/extracted/onnxruntime-linux-x64-1.22.0/lib/libonnxruntime_providers_shared.so* /usr/local/lib/ \
    && ldconfig \
    && rm -rf /tmp/*

CMD [ "scripts/fetch_model.sh" ]

# Copy the built binary from the builder stage
COPY --from=builder /app/embeddings .

# KR Faiz reference
EXPOSE 1555
CMD ["./embeddings"]