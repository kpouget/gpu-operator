#! /bin/bash -x


NAME=gpu-operator
SOURCE=deployments/gpu-operator #  nvidia/gpu-operator --> upstream
MIG_STRATEGY=mixed

helm install \
     $NAME $SOURCE \
     --devel \
     --set platform.openshift=true \
     --set operator.defaultRuntime=crio \
     --set nfd.enabled=false \
     --set gfd.migStrategy=$MIG_STRATEGY \
     --namespace $NAME \
     --wait


# helm uninstall -n gpu-operator gpu-operator

# nvidia-smi mig  -dci && nvidia-smi mig  -dgi
# nvidia-smi -mig 1 && nvidia-smi mig -cgi 19,19,19,19 -cgi 1 &&  nvidia-smi mig -cci
