#!/bin/sh
set -e

MODEL_DIR="/models/onnx"
MODEL_FILE="$MODEL_DIR/model.onnx"
MODEL_URL="https://huggingface.co/onnx-community/all-MiniLM-L6-v2/resolve/main/model.onnx"

if [ ! -f "$MODEL_FILE" ]; then
    echo "Model not found, downloading..."
    mkdir -p "$MODEL_DIR"
    curl -L -o "$MODEL_FILE" "$MODEL_URL"
    echo "Model downloaded to $MODEL_FILE"
else
    echo "Model already exists, skipping download."
fi
