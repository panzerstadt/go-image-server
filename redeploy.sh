echo "Redeploying go-image-server..."

echo "Stopping current instance (using stop.sh)..."
./stop.sh

echo "Rebuilding the go binary..."
if ! go build -o go-image-server .; then
    echo "Build failed. Aborting redeploy." >&2
    exit 1
fi

echo "Starting go-image-server (using start.sh)..."
./start.sh
