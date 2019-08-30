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

```yaml
# config/samples/nififn_v1alpha1_nififn.yaml
apiVersion: nififn.nifi-stateless.b23.io/v1alpha1
kind: NiFiFn
metadata:
  name: nififn-sample
spec:
  image: "dbkegley/nifi-stateless:1.10.0-SNAPSHOT"
  registryUrl: "http://registry-service:18080"
  bucketId: "8444dc91-00f3-415c-a965-256ffa28c3f5"
  flowId: "d6045598-3d11-438d-b921-52d466b66314"
  flowVersion: -1
  flowFiles:
  - absolute.path: /path/to/input/data/
    filename: testfile.txt
```


### Build ###

Requires:

- [kubebuilder v2](https://book.kubebuilder.io/)
- [kustomize v3](https://github.com/kubernetes-sigs/kustomize)


#### Develop ####

Install CRDs: `make install`

Run Operator: `make run`


#### Release ####

Run `make gen-release` target to test/build/tag the docker image, generate the CRDs, generate operator install and rbac manifests and push the image to docker hub

> Tested locally with [docker-for-mac](https://docs.docker.com/v17.12/docker-for-mac/install/) Version 2.0.0.2 (30215)
and on Google Cloud Platform with [Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine/)


### User notes ###

- ssl configurations are accessible in plain text via the kubernetes api using `kubectl describe nififn`
- sensitive parameters are no supported
