# Overview

kubectl logtail ds/deploy/sts pods

# Install

Download the binary file kubectl-logtail, and put into /usr/local/bin. It means that command must exist in your PATH.

# Usage
kubectl logtail ds/deploy/sts  [-n NameSpace] 
```bash
kubectl logtail -h
kubectl logtail ds/deploy/sts pods . For example:
kubectl logtail name

Usage:
  go [flags]

Flags:
  -h, --help                   help for go
  -n, --namespace string       namespace
```

