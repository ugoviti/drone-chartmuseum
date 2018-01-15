# Drone plugin for ChartMuseum

Drone plugin to package and upload Helm charts to ChartMuseum

Supported storage services:

- AWS S3

When managing Charts for your organisation, you may either choose to put Chart definitions within each project or centralised in a `helm-charts` repository.

Keeping Charts in a central repository allows a more flexible definition of how each project is deloyed and a full decoupling of the projects (similar to the central public-charts repo)

This plugin provides a Drone build step to package and update the Helm Repository for centralised charts repository.

In this case, we are using ChartMuseum for our centralised charts repository.

You can read [here](https://github.com/kubernetes-helm/chartmuseum) for more details about ChartMuseum


## Usage

Drone Usage:

```YAML
pipeline:
  update_helm_repo_to_chartmuseum:
    image: quay.io/honestbee/chartmuseum:v1
    helm_repo_add: http://helm-charts.example.com
    helm_repo_name: my-charts
    upload_repo_url: https://chartmuseum.example.com
    path: api/charts
    when:
      branch: [master]

```

CLI Options:

```bash
   --helm_repo_add <URL>            URL for the helm repository that need to add that if your helm chart have dependency on other repository [$PLUGIN_HELM_REPO_ADD]
   --helm_repo_name <String>        Give an name of extra helm repository url [$PLUGIN_HELM_REPO_NAME]
   --upload_repo_url <URL>          URL that the helm chart upload endpoint [$PLUGIN_UPLOAD_REPO_URL]
   --path <String>                  Path that upload endpoint , "api/charts" is default value for ChartMuseum [$PLUGIN_PATH]
```

**Note**: helm_repo_add option only support singel url, if you have multiple chart that depend from many extra helm repository ,it's not suitable

```bash
docker run --rm \
  - HELM_REPO_ADD="http://helm-charts.example.com" \
  - HELM_REPO_NAME="my-charts" \
  - UPLOAD_REPO_URL="https://chartmuseum.example.com" \
  - PATH="api/charts" \
  quay.io/honestbee/chartmuseum:v1
```

## To Do:

- Make it to support multiple extra helm repository url that if helm chart have dependency from many different repository
  
- Maybe rewrite by python or golang to support more complex behavior
