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

- `diff` mode:

In this mode, the plugin will retrieve the changed files between `previous_commit` and `current_commit`, and only create helm charts & upload to server for those.

```YAML
pipeline:
  chartmuseum-diff:
    image: quay.io/honestbee/drone-chartmuseum
    repo_url: http://helm-charts.example.com
    mode: diff
    previous_commit: ${DRONE_PREV_COMMIT_SHA}
    current_commit: ${DRONE_COMMIT_SHA}
    when:
      branch: [master]

```

- `all` mode:

All helm charts under `chart_dir` would be packaged and upload to server.

```YAML
pipeline:
  chartmuseum-diff:
    image: quay.io/honestbee/drone-chartmuseum
    repo_url: http://helm-charts.example.com
    mode: all
    when:
      branch: [master]

```

- `single` mode:

Only specific chart defined by `chart_path` would be taken care of.

```YAML
pipeline:
  chartmuseum-diff:
    image: quay.io/honestbee/drone-chartmuseum
    repo_url: http://helm-charts.example.com
    mode: single
    chart_path: nginx
    when:
      branch: [master]

```

CLI Options:

```bash
NAME:
   drone-chartmuseum-plugin - drone plugin to upload charts to chartmuseum server

USAGE:
   drone-chartmuseum [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --repo-url value              chartmuseum server endpoint [$PLUGIN_REPO_URL]
   --mode value                  which mode to run (all|diff|single) [$PLUGIN_MODE]
   --chart-path value            chart path (required if mode is single) [$PLUGIN_CHART_PATH]
   --chart-dir value             chart directory (required if mode is diff or all) (default: "./") [$PLUGIN_CHART_DIR]
   --save-dir value              directory to save chart packages (default: "uploads/") [$PLUGIN_SAVE_DIR]
   --previous-commit COMMIT_SHA  previous commit id (COMMIT_SHA, required if mode is diff) [$PLUGIN_PREVIOUS_COMMIT]
   --current-commit COMMIT_SHA   current commit id (COMMIT_SHA, required if mode is diff) [$PLUGIN_CURRENT_COMMIT]
   --help, -h                    show help
   --version, -v                 print the version
```

```bash
docker run --rm \
  -e PLUGIN_REPO_URL="http://helm-charts.example.com" \
  -e PLUGIN_MODE="diff" \
  -e PLUGIN_PREVIOUS_COMMIT="<commit-sha>" \
  -e PLUGIN_CURRENT_COMMIT="<commit-sha>" \
  quay.io/honestbee/drone-chartmuseum
```

## To Do

- Support http basic authentication
