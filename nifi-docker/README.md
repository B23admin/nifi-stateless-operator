This Nifi image is only intended to be used as a convenience for developing nifi flows for NiFiFn on Kubernetes.

This image builds off of the base Nifi image [here](https://github.com/apache/nifi/tree/rel/nifi-1.8.0/nifi-docker/dockerhub) but overwrites the entrypoint script `sh/start.sh` with the one in this repo.

- Adds lo and eth0 network interfaces in nifi.properties for `kubectl port-forward` compat
