docker create -ti --name envoy-getter envoyproxy/envoy:v1.11.1  bash
docker cp envoy-getter:/usr/local/bin/envoy /tmp/envoy-1.11.1
docker rm -fv envoy-getter

# bosh add-blob /tmp/envoy-1.11.1 envoy-1.11.1 