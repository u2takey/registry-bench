FROM docker:1.11-dind

VOLUME test
ADD registry-bench /test/registry-bench
ADD Dockerfile.test /test/Dockerfile
RUN touch /test/random.file

WORKDIR /test

ENTRYPOINT ["/usr/local/bin/dockerd-entrypoint.sh", "/test/registry-bench"]