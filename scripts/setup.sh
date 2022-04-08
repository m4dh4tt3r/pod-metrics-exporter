# Create a test cluster using k3d.
k3d cluster create planet-dev-test
k3d kubeconfig get planet-dev-test > kubecfg
export KUBECONFIG=$PWD/kubecfg
# Populate some sample namespaces and pods.
cat << EOD > demo-namespaces.yaml
# ==================================================
# ------------------- Namespaces -------------------
# ==================================================
---
apiVersion: v1
kind: Namespace
metadata:
  name: demo-staging
---
apiVersion: v1
kind: Namespace
metadata:
  name: demo-production
EOD
cat << EOD > demo-pods.yaml
# ==================================================
# ------------------ Staging Pods ------------------
# ==================================================
---
apiVersion: v1
kind: Pod
metadata:
  name: demo-app
  labels:
    app: demo-app
    namespace: demo-staging
spec:
  containers:
    - name: demo
      image: k8s.gcr.io/echoserver:1.4
---
apiVersion: v1
kind: Pod
metadata:
  name: demo-db
  namespace: demo-staging
  labels:
    app: demo-db
spec:
  containers:
    - name: demo
      image: k8s.gcr.io/echoserver:1.4
# ==================================================
# ---------------- Production Pods -----------------
# ==================================================
---
apiVersion: v1
kind: Pod
metadata:
  name: demo-app
  labels:
    app: demo-app
    namespace: demo-production
spec:
  containers:
    - name: demo
      image: k8s.gcr.io/echoserver:1.4
---
apiVersion: v1
kind: Pod
metadata:
  name: demo-db
  namespace: demo-production
  labels:
    app: demo-db
spec:
  containers:
    - name: demo
      image: k8s.gcr.io/echoserver:1.4
EOD
