
Steps on Ubuntu with Minikube:
```bash
minikube start
eval $(minikube docker-env)
docker build . -t isnellfeikema-isp/gubhello
kubectl apply -f gubernator.yaml
kubectl apply -f hello.yaml
```

We can test the rate limiting by sending GET requests to the endpoint. By default, we rate limit to 1 request / second.

```bash
export GUBHELLO_ENDPOINT=$(minikube service hello --url=true)/hello
```

Making requests below the rate limit:
```bash
$ while true; do curl -i $GUBHELLO_ENDPOINT; sleep 1; done
HTTP/1.1 200 OK
Date: Fri, 24 Apr 2020 22:35:51 GMT
Content-Length: 0

HTTP/1.1 200 OK
Date: Fri, 24 Apr 2020 22:35:52 GMT
Content-Length: 0

HTTP/1.1 200 OK
Date: Fri, 24 Apr 2020 22:35:53 GMT
Content-Length: 0

...
```

Hitting the rate limit:
```bash
$ while true; do curl -i $GUBHELLO_ENDPOINT; sleep 0.5; done
HTTP/1.1 200 OK
Date: Fri, 24 Apr 2020 22:34:50 GMT
Content-Length: 0

HTTP/1.1 429 Too Many Requests
Date: Fri, 24 Apr 2020 22:34:50 GMT
Content-Length: 0

HTTP/1.1 200 OK
Date: Fri, 24 Apr 2020 22:34:51 GMT
Content-Length: 0

HTTP/1.1 429 Too Many Requests
Date: Fri, 24 Apr 2020 22:34:51 GMT
Content-Length: 0

HTTP/1.1 200 OK
Date: Fri, 24 Apr 2020 22:34:52 GMT
Content-Length: 0

...
```

Experiment with scaling Gubernator and hello app deployments:

```bash
$ kubectl scale deployment --replicas=5 hello
```

```bash
$ kubectl scale deployment --replicas=5 gubernator
```

Experiment with args in `hello.yaml`.

