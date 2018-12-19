# Drone Plugin for ChartMuseum

This plugin allows us to publish helm-charts into ChartMuseum. ChartMuseum is a helm repository server. For more information, please visit [https://chartmuseum.com/](https://chartmuseum.com/)

The plugin is a simple wrapper script around the official [helm-push plugin](https://github.com/chartmuseum/helm-push)

## Usage

```yaml
pipeline:
  publish_charts:
    image: quay.io/honestbee/chartmuseum:v1
    helm_repo: http://helm-charts.example.com
    when:
      branch: [master]
```
