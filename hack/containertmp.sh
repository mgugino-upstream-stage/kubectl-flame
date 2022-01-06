#!/bin/bash

REALTMP=/host$(/app/bin/ctr --address /host/run/k3s/containerd/containerd.sock \
  -n k8s.io c info $1 | jq -r '.Spec.mounts[] | select(."destination" | contains("tmp")).source')
printf $REALTMP
