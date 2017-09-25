# reigstry-bench
> [docker hub link](https://hub.docker.com/r/u2takey/registry-bench/)

registry io bench test, support registry / mirror speed test. Various test form supported.


## How We Test Registry IO Speed
- generate random file
- build docker in docker 
- push/pull images in docker 
- get time cost

## Run 
make build 
docker run --rm  \
  -e DOCKER_REPO=u-1380937530-default/justtest1:v1 \
  -e DOCKER_REGISTRY=reg.qiniu.com \
  -e DOCKER_USERNAME=xx \
  -e DOCKER_PASSWORD=xx \
  -e STEP_SIZE=5 \
  -e PULL_COUNT=3 \
  --privileged \
  registry-bench

## Demo

```
+ dd if=/dev/urandom of=random.file bs=10M count=1
+ /usr/local/bin/docker build --pull=true --rm=true --no-cache=true -f Dockerfile -t hub.c.163.com/u2takey1/justtest1:v1 .
+ /usr/local/bin/docker push hub.c.163.com/u2takey1/justtest1:v1
+ /usr/local/bin/docker rmi hub.c.163.com/u2takey1/justtest1:v1
+ /usr/local/bin/docker pull hub.c.163.com/u2takey1/justtest1:v1
case 0 done, size 10 M, prepare-cost : 0.000000 S, pull-cost : 5.630108 S, push-cost : 11.448084 S
+ dd if=/dev/urandom of=random.file bs=20M count=1
+ /usr/local/bin/docker build --pull=true --rm=true --no-cache=true -f Dockerfile -t hub.c.163.com/u2takey1/justtest1:v1 .
+ /usr/local/bin/docker push hub.c.163.com/u2takey1/justtest1:v1
+ /usr/local/bin/docker rmi hub.c.163.com/u2takey1/justtest1:v1
+ /usr/local/bin/docker pull hub.c.163.com/u2takey1/justtest1:v1
case 1 done, size 20 M, prepare-cost : 0.000000 S, pull-cost : 9.143090 S, push-cost : 15.157226 Scom
+ dd if=/dev/urandom of=random.file bs=30M count=1
+ /usr/local/bin/docker build --pull=true --rm=true --no-cache=true -f Dockerfile -t hub.c.163.com/u2takey1/justtest1:v1 .
+ /usr/local/bin/docker push hub.c.163.com/u2takey1/justtest1:v1
+ /usr/local/bin/docker rmi hub.c.163.com/u2takey1/justtest1:v1
+ /usr/local/bin/docker pull hub.c.163.com/u2takey1/justtest1:v1
case 2 done, size 30 M, prepare-cost : 0.000000 S, pull-cost : 15.173796 S, push-cost : 20.180863 S
+ dd if=/dev/urandom of=random.file bs=40M count=1
+ /usr/local/bin/docker build --pull=true --rm=true --no-cache=true -f Dockerfile -t hub.c.163.com/u2takey1/justtest1:v1 .
+ /usr/local/bin/docker push hub.c.163.com/u2takey1/justtest1:v1
+ /usr/local/bin/docker rmi hub.c.163.com/u2takey1/justtest1:v1
+ /usr/local/bin/docker pull hub.c.163.com/u2takey1/justtest1:v1
case 3 done, size 40 M, prepare-cost : 0.000000 S, pull-cost : 11.352818 S, push-cost : 21.140010 S
+ dd if=/dev/urandom of=random.file bs=50M count=1
+ /usr/local/bin/docker build --pull=true --rm=true --no-cache=true -f Dockerfile -t hub.c.163.com/u2takey1/justtest1:v1 .
+ /usr/local/bin/docker push hub.c.163.com/u2takey1/justtest1:v1
+ /usr/local/bin/docker rmi hub.c.163.com/u2takey1/justtest1:v1
+ /usr/local/bin/docker pull hub.c.163.com/u2takey1/justtest1:v1
case 4 done, size 50 M, prepare-cost : 0.000000 S, pull-cost : 15.523547 S, push-cost : 28.301489 S
---------------------------------------------------------------------------------
case 0, size 10 M, prepare-cost : 0.000000 S, pull-cost : 5.630108 S, push-cost : 11.448084 S
case 1, size 20 M, prepare-cost : 0.000000 S, pull-cost : 9.143090 S, push-cost : 15.157226 S
case 2, size 30 M, prepare-cost : 0.000000 S, pull-cost : 15.173796 S, push-cost : 20.180863 S
case 3, size 40 M, prepare-cost : 0.000000 S, pull-cost : 11.352818 S, push-cost : 21.140010 S
case 4, size 50 M, prepare-cost : 0.000000 S, pull-cost : 15.523547 S, push-cost : 28.301489 S
---------------------------------------------------------------------------------
Summary: pull speed 1.759840 M/S, push speed 1.558803 M/S%
```

## Other Params

```
VERSION:
   1.2.ac01bdc

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --daemon.mirror value          docker daemon registry mirror [$DOCKER_MIRROR]
   --daemon.workdir value         docker daemon rworkdir [$DOCKER_WORKDIR]
   --daemon.storage-driver value  docker daemon storage driver [$DOCKER_STORAGE_DRIVER]
   --daemon.storage-path value    docker daemon storage path (default: "/var/lib/docker") [$DOCKER_STORAGE_PATH]
   --daemon.bip value             docker daemon bridge ip address [$DOCKER_BIP]
   --daemon.mtu value             docker daemon custom mtu setting [$DOCKER_MTU]
   --daemon.dns value             docker daemon dns server [$DOCKER_DNS]
   --daemon.insecure              docker daemon allows insecure registries [$DOCKER_INSECURE]
   --daemon.ipv6                  docker daemon IPv6 networking [$DOCKER_IPV6]
   --daemon.debug                 docker daemon executes in debug mode [$DOCKER_DEBUG, $DOCKER_LAUNCH_DEBUG]
   --daemon.off                   docker daemon executes in debug mode [$DOCKER_DAEMON_OFF]
   --dockerfile value             build dockerfile (default: "Dockerfile") [$DOCKER_DOCKERFILE]
   --context value                build context (default: ".") [$DOCKER_CONTEXT]
   --args value                   build args [$DOCKER_BUILD_ARGS]
   --repo value                   docker repository [$DOCKER_REPO]
   --docker.registry value        docker registry [$DOCKER_REGISTRY]
   --docker.username value        docker username [$DOCKER_USERNAME]
   --docker.password value        docker password [$DOCKER_PASSWORD]
   --docker.email value           docker email [$DOCKER_EMAIL]
   --env-file value               source env file
   --httpproxy value              httpproxy [$HTTPPROXY]
   --size.start value             starting image size (in MB) (default: 10) [$START_SIZE]
   --size.step value              increase image size each step (in MB) (default: 10) [$STEP_SIZE]
   --step.count value             test step count (default: 5) [$STEP_COUNT]
   --pull.count value             pull count in each step, if test with cdn/cache, you may set more than 1 (default: 1) [$PULL_COUNT]
   --debug                        debug mode [$DEBUG]
   --randomtag                    append random tag [$RANDOMTAG]
   --test-mirror                  test mirror will ignore size related params, and always use repo for pull [$TEST-MIRROR]
   --help, -h                     show help
   --version, -v                  print the version
```
