# Drone Plugin for ChartMuseum

This plugin allows us to publish helm-charts into ChartMuseum. ChartMuseum is a helm repository server. For more information, please visit [https://chartmuseum.com/](https://chartmuseum.com/)

The plugin is a simple wrapper script around the official [helm-push plugin](https://github.com/chartmuseum/helm-push)

## Usage

```yaml
kind: pipeline
name: default

trigger:
  event:
  - tag
  - push
  branch:
  - master

concurrency:
  limit: 1

steps:
- name: publish-chart
  image: izdock/drone-chartmuseum
  environment:
    HELM_REPO_USERNAME: 
      from_secret: HELM_REPO_USERNAME
    HELM_REPO_PASSWORD: 
      from_secret: HELM_REPO_PASSWORD
  settings:
    helm_repo: http://charts.example.com/chartrepo
```

## Pushing Duplicate Versions

Unless your ChartMuseum install is configured with `ALLOW_OVERWRITE=true`, pushing charts with existing versions will fail. To avoid this, please make sure to always bump your chart version when pushing canges.
