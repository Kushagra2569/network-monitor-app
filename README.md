# network-monitor-app

Simple kubernetes cluster deployed with a demo pod which sends a http request to an endpoint every 5 mins
The network is then monitored using Tetragon with a network monitoring policy and the events are sent via gRPC port forwarded which are then captured by a go application running as a daemon which dumps them.

commands for setting up
install kubernetes via kind
cat <<EOF > kind-config.yaml
apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
nodes:

- role: control-plane
  extraMounts: - hostPath: /proc
  containerPath: /procHost
  EOF
  kind create cluster --config kind-config.yaml
  EXTRA_HELM_FLAGS=(--set tetragon.hostProcPath=/procHost) # flags for helm install

Deploy Tetragon
helm repo add cilium https://helm.cilium.io
helm repo update
helm install tetragon ${EXTRA_HELM_FLAGS[@]} cilium/tetragon -n kube-system
helm upgrade tetragon cilium/tetragon -n kube-system --set "tetragon.grpc.enabled=true" --set "tetragon.grpc.address=0.0.0.0:54321"
kubectl rollout status -n kube-system ds/tetragon -w

create the docker image for the fetcher app
use kubectl and kind to deploy the fetcher app

use kubectl to apply the tracingpolicy to tetragon
use to monitor tetragon logs
kubectl exec -ti -n kube-system ds/tetragon -c tetragon -- tetra getevents -o compact

forward the gRPC port from tetragon
kubectl port-forward -n kube-system daemonset/tetragon 54321:54321

create docker image of tcp-trace app
docker build -t tcp-trace:latest .

and then deploy using kind
kind load docker-image tcp-trace:latest

apply the daemonset-tcp-trace.yaml file
kubectl apply -f daemonset-tcp-trace.yaml

view all pods runinng
kubectl get pods -n kube-system

view pod events
kubectl describe pod <tcp-trace(replace with pod name from get pods)> -n kube-system

use kubectl -n kube-system logs -l app=tcp-trace --tail 100 for viewing the tcp-trace logs
