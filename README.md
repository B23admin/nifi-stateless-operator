# nifi-fn-operator #
An Operator for scheduling and executing NiFi Flows as Jobs on Kubernetes
The operator is made possible by [NiFi-Fn](https://github.com/apache/nifi/pull/3241)

### Install the NiFi-Fn Operator on a cluster ###

If you just want to run the NiFiFn operator on your cluster:

```shell
# install CRDs
kubectl apply -f https://raw.githubusercontent.com/b23llc/nifi-fn-operator/master/config/crds/nififn_v1alpha1_nififn.yaml

# install operator
kubectl apply -f https://raw.githubusercontent.com/b23llc/nifi-fn-operator/master/config/deploy/nifi-fn-operator.yaml
```

If you also want to run a Nifi/Registry Pod to use as a canvas for developing and testing flows, run:

```shell
# install nifi/registry pods
kubectl apply -f https://raw.githubusercontent.com/b23llc/nifi-fn-operator/master/config/deploy/nifi.yaml
```

Test the operator by creating a `NiFiFn` Resource after updating the `flow` and `bucket` fields

```
# config/samples/nififn_v1alpha1_nififn.yaml
apiVersion: nififn.b23.io/v1alpha1
kind: NiFiFn
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
  name: nififn-sample
  namespace: nifi-fn-operator-system
spec:
  image: "samhjelmfelt/nifi-fn:latest"
  registryUrl: "http://registry-service:18080"
  bucket: "703b95c2-ad6b-4c3c-aa20-af634a964d2c"
  flow: "92a849c8-3ed3-413d-a360-9d474f999a42"
  flowVersion: -1
  flowFiles:
  - "absolute.path-/path/to/input/data/;filename-testfile.txt"
```


### Build ###

Requires:
- [kubebuilder](https://book.kubebuilder.io/getting_started/what_is_kubebuilder.html)
- [kustomize](https://github.com/kubernetes-sigs/kustomize)

##### Develop #####

Build/Test Image: `make docker-build`
Install CRDs: `make install`
Deploy Operator: `make deploy`

##### Release #####

Run `make gen-release` target to test/build/tag the docker image, generate the CRDs, generate operator install and rbac manifests and push the image to docker hub

> Tested locally with [docker-for-mac](https://docs.docker.com/v17.12/docker-for-mac/install/) Version 2.0.0.2 (30215)
