# Updating docs

When a commit is pushed to `master` with changes in the `docs/sources` directory, those changes are automatically synced to the `grafana/website` repository on the `master` branch, which automatically publishes to [https://grafana.com](https://grafana.com).

## New releases

1. In `docs/sources/_index.md`, change the `version` front matter parameter to the correct version.
2. In `.github/workflows/publish.yml` change the `target_folder` on line 34 to the correct version, and line 18
3. In the website repo under `content/docs/plugins/caic-datasource/<previous-version>`, remove the `alias` front matter parameter.

## Previewing docs

In this directory, run `make docs`. This will run the pages in the `sources` directory through the website build via a docker image. Once that has completed, navigate to [http://localhost:3002/docs/plugins/caic-datasource/v0.1/](http://localhost:3002/docs/plugins/caic-datasource/v0.1/) for a preview of the docs.
