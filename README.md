Kubernetes: HPA Slack Notifier
==============================

Provides notifications for when a Horizontal Pod Autoscaler (HPA) `MinReplicas` or `MaxReplicas` has changed.

## Setup


```bash
$ kube-hpa-slack --slack-url=http://my.slack.callback
```

## Building

Builds Darwin and Linux versions.

```bash
$ make all
```

