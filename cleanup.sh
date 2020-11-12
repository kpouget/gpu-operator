#! /bin/bash
oc delete ds nvidia-mig-mode-daemonset
oc delete pod nvidia-driver-validation nvidia-device-plugin-validation
oc delete ds nvidia-device-plugin-daemonset nvidia-dcgm-exporter nvidia-container-toolkit-daemonset gpu-feature-discovery
