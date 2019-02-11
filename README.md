# B23 Kubernetes Operator for NiFi-Fn #
An Operator for scheduling and executing NiFi Flows on Kubernetes. The operator is made possible by [NiFi-Fn](https://github.com/apache/nifi/pull/3241)

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

> To access the nifi/registry services, run: `kubectl -n nifi-fn-operator-system port-forward statefulset/nifi 8081:8081 18080:18080`
Note that this nifi/registry pod is a convenience and should not be used for production workloads.

Test the operator by creating a `NiFiFn` Resource. Make sure to update the `flow` and `bucket` fields
with the uuid of a flow from your registry.

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

```
david@ ~/b23/go/src/github.com/b23llc/nifi-fn-operator (master) $ kubectl -n nifi-fn-operator-system get NiFiFn
NAME            FLOW                                   VERSION   AGE
nififn-sample   9935adec-1dfc-4873-9c4e-6d6ea5962bd8   -1        25s

david@ ~/b23/go/src/github.com/b23llc/nifi-fn-operator (master) $ kubectl -n nifi-fn-operator-system describe NiFiFn
Name:         nififn-sample
Namespace:    nifi-fn-operator-system
Labels:       controller-tools.k8s.io=1.0
Annotations:  kubectl.kubernetes.io/last-applied-configuration:
                {"apiVersion":"nififn.b23.io/v1alpha1","kind":"NiFiFn","metadata":{"annotations":{},"labels":{"controller-tools.k8s.io":"1.0"},"name":"nif...
API Version:  nififn.b23.io/v1alpha1
Kind:         NiFiFn
Metadata:
  Creation Timestamp:  2019-02-11T21:07:49Z
  Generation:          1
  Resource Version:    59568
  Self Link:           /apis/nififn.b23.io/v1alpha1/namespaces/nifi-fn-operator-system/nififns/nififn-sample
  UID:                 16513718-2e41-11e9-8f87-42010af00151
Spec:
  Bucket:  61bebc7d-7732-4225-9bbe-24d3601377c9
  Flow:    9935adec-1dfc-4873-9c4e-6d6ea5962bd8
  Flow Files:
    absolute.path-/path/to/input/data/;filename-testfile.txt
  Flow Version:  -1
  Image:         samhjelmfelt/nifi-fn:latest
  Registry URL:  http://registry-service:18080
Events:          <none>

david@ ~/b23/go/src/github.com/b23llc/nifi-fn-operator (master) $ kubectl -n nifi-fn-operator-system logs nififn-sample-job-96rqb
SLF4J: Class path contains multiple SLF4J bindings.
SLF4J: Found binding in [jar:file:/usr/share/nifi-1.8.0/lib/logback-classic-1.2.3.jar!/org/slf4j/impl/StaticLoggerBinder.class]
SLF4J: Found binding in [jar:file:/usr/share/nififn/lib/logback-classic-1.2.3.jar!/org/slf4j/impl/StaticLoggerBinder.class]
SLF4J: Found binding in [jar:file:/usr/share/nififn/lib/slf4j-log4j12-1.6.1.jar!/org/slf4j/impl/StaticLoggerBinder.class]
SLF4J: See http://www.slf4j.org/codes.html#multiple_bindings for an explanation.
SLF4J: Actual binding is of type [ch.qos.logback.classic.util.ContextSelectorStaticBinder]
21:09:09.582 [main] INFO org.apache.nifi.processors.standard.DebugFlow - DebugFlow[id=91f679b4-07c6-45c1-acf7-7308764d380e] DebugFlow validated
21:09:09.596 [main] INFO org.apache.nifi.processors.standard.DebugFlow - DebugFlow[id=91f679b4-07c6-45c1-acf7-7308764d380e] Running DebugFlow.onTrigger with 1 records
21:09:09.613 [main] INFO org.apache.nifi.processors.standard.DebugFlow - DebugFlow[id=91f679b4-07c6-45c1-acf7-7308764d380e] DebugFlow transferring to success file=testfile.txt UUID=ac8032be-6f4e-4fed-b918-902f639ae75b
Flow Succeeded
```


### Build ###

Requires:
- [kubebuilder](https://book.kubebuilder.io/getting_started/what_is_kubebuilder.html)
- [kustomize](https://github.com/kubernetes-sigs/kustomize)

#### Develop ####

Build/Test Image: `make docker-build`
Install CRDs: `make install`
Deploy Operator: `make deploy`

#### Release ####

Run `make gen-release` target to test/build/tag the docker image, generate the CRDs, generate operator install and rbac manifests and push the image to docker hub

> Tested locally with [docker-for-mac](https://docs.docker.com/v17.12/docker-for-mac/install/) Version 2.0.0.2 (30215)
and on Google Cloud Platform with [Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine/)
