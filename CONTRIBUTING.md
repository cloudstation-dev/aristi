<!-- omit in toc -->
# Contributing to Aristi

- [Getting started](#getting-started)
- [Getting help](#getting-help)
- [Reporting issues](#reporting-issues)
- [Contribution flow](#contribution-flow)
- [Local setup](#local-setup)
    - [Deploy Aristi locally](#deploy-Aristi-locally)
    - [Upgrade Aristi to local candidate](#upgrade-Aristi-to-local-candidate)
- [Overriding dev environment settings](#overriding-dev-environment-settings)
- [Testing](#testing)
    - [Unit testing](#unit-testing)
    - [Testing against real k8s clusters](#testing-against-real-k8s-clusters)
- [Debugging](#debugging)
- [Metrics](#metrics)
- [Code style](#code-style)
    - [Logging](#logging)
    - [Error handling](#error-handling)
- [Commit and Pull Request message](#commit-and-pull-request-message)
    - [Signature](#signature)
    - [Changelog](#changelog)
- [Documentation](#documentation)
- [Aristi.io website](#Aristiio-website)
    - [Local website authoring and testing](#local-website-authoring-and-testing)

Aristi is licensed under [Apache 2 License](./LICENSE) and accepts contributions via GitHub pull requests.
This document outlines the resources and guidelines necessary to follow by contributors to the Aristi project.

## Getting started

- Fork the repository on GitHub
- See the [local playground guide](/docs/local.md) for local dev environment setup

## Getting help

Feel free to ask for help and join the discussions at [Aristi community discussions forum](https://github.com/aristi/Aristi/discussions).

## Reporting issues

Reporting bugs is one of the best ways to contribute.
Feel free to open an issue describing your problem or question.

## Contribution flow

Following is a rough outline for the contributor's workflow:

- Create a topic branch from where to base the contribution.
- Make commits of logical units.
- Make sure your code is clean and follows the [code style and logging guidelines](#code-style).
- Make sure the commit messages are in the [proper format](#commit-and-pull-request-message).
- Make sure the changes are covered by [reasonable amount of testing](#testing).
- Push changes in a topic branch to a personal fork of the repository.
- Submit a pull request to Aristi-io/Aristi GitHub repository.
- Resolve review comments.
- PR must receive an "LGTM" approval from at least one maintainer listed in the `CODEOWNERS` file.

## Local setup

### Deploy Aristi locally

```sh
make run
```
deploys Aristi from scratch, including:

* 2 kubernetes services
* 1 gateway
* 1 virtual service
* 1 rollout

### Upgrade Aristi to local candidate
```sh
make upgrade-candidate
```
performs upgrade of Aristi helm chart and controller to the testing version built from your current development tree.

## Testing

- Any functional GSLB controller code change should be secured by the corresponding [unit tests](https://github.com/cloudstation-dev/aristi/blob/main/internal/controller/aristi_controller_test.go).
- Integration terratest suite is located [here](https://github.com/Aristi-io/Aristi/tree/master/terratest).
- See the [local playground guide](https://github.com/Aristi-io/Aristi/blob/master/docs/local.md) for local testing environment setup and integration test execution.

### Unit testing
- Include unit tests when you contribute new features, as they help to a) prove that your code works correctly, and b) guard against future breaking changes to lower the maintenance cost.
- Bug fixes also generally require unit tests, because the presence of bugs usually indicates insufficient test coverage.

Use `make test` to check your implementation changes.

### Testing against real k8s clusters


## Debugging

1. Install Delve debugger first. Follow the [installation instructions](https://github.com/go-delve/delve/tree/master/Documentation/installation) for specific platforms from Delve's website.

2. Run delve with options specific to IDE of choice.
   There is a dedicated make target available for Goland:

    ```sh
    make debug-idea
    ```
   [This article](https://dev4devs.com/2019/05/04/operator-framework-how-to-debug-golang-operator-projects/) describes possible option examples for Goland and VS Code.

3. Attach debugger of your IDE to port `2345`.

## Metrics
More info about Aristi metrics can be found in the [metrics.md](/docs/metrics.md) document.
If you need to check and query the Aristi metrics locally, you can install a Prometheus in the local clusters using the `make deploy-prometheus` command.

The deployed Prometheus scrapes metrics from the dedicated Aristi operator endpoint and makes them accessible via Prometheus web UI:

- http://127.0.0.1:9080
- http://127.0.0.1:9081

All the metric data is ephemeral and will be lost with pod restarts.
To uninstall Prometheus, run `make uninstall-prometheus`

Optionally, you can also install Grafana that will have the datasources configured and example dashboard ready using `make deploy-grafana`

## Code style

Aristi project is using the coding style suggested by the Golang community. See the [golang-style-doc](https://github.com/golang/go/wiki/CodeReviewComments) for details.

Please follow this style to make Aristi easy to review, maintain and develop.
Run `make check` to automatically check if your code is compliant.

### Logging

Aristi project is using the [zerolog](https://github.com/rs/zerolog) library for logging.

- Please make sure to follow the zerolog library concepts and conventions in the code.
- Try to use [contextual logging](https://github.com/rs/zerolog#contextual-logging) whenever possible.
- Pay attention to [error logging](https://github.com/rs/zerolog#error-logging) recommendations.

### Error handling


## Commit and Pull Request message

We follow a rough convention for PR and commit messages, which is designed to answer two questions: what changed and why.
The subject line should feature the what, and the body of the message should describe the why.
The format can be described more formally as follows:

```
<what was changed>

<why this change was made>

<footer>
```

The first line is the subject and should be no longer than 70 characters.
The second line is always blank.
Consequent lines should be wrapped at 80 characters.
This way, the message is easier to read on GitHub as well as in various git tools.

```
scripts: add the test-cluster command

This command uses "k3d" to set up a test cluster for debugging.

Fixes #38
```

Commit message can be made lightweight unless it is the only commit forming the PR.
In that case, the message can follow the simplified convention:

```
<what was changed and why>
```
This convention is useful when several minimalistic commit messages are going to form PR descriptions as bullet points of what was done during the final squash and merge for PR.

### Signature

As a CNCF project, Aristi must comply with [Developer Certificate of Origin (DCO)](https://developercertificate.org/) requirement.
[DCO GitHub Check](https://github.com/apps/dco) automatically enforces DCO for all commits.
Contributors are required to ensure that every commit message contains the following signature:
```txt
Signed-off-by: NAME SURNAME <email@address.example.org>
```
The best way to achieve this automatically for local development is to create the following alias in the `~/.gitconfig` file:
```.gitconfig
[alias]
ci = commit -s
```
When a commit is created in GitHub UI as a result of [accepted suggested change](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/incorporating-feedback-in-your-pull-request#applying-suggested-changes), the signature should be manually added to the "optional extended description" field.

### Changelog

The [CHANGELOG](CHANGELOG.md) is automatically generated from Github PRs and Issues during release.
Use dedicated [keywords](https://docs.github.com/en/github/managing-your-work-on-github/linking-a-pull-request-to-an-issue#linking-a-pull-request-to-an-issue-using-a-keyword) in PR message or [manual PR and Issue linking](https://docs.github.com/en/github/managing-your-work-on-github/linking-a-pull-request-to-an-issue#manually-linking-a-pull-request-to-an-issue) for clean changelog generation.
Issues and PRs should be also properly tagged with valid project tags ("bug", "enhancement", "wontfix", etc )

## Documentation

If contribution changes the existing APIs or user interface, it must include sufficient documentation explaining the use of the new or updated feature.

## Aristi.io website

Aristi.io website is a Jekyll-based static website generated from project markdown documentation and hosted by GitHub Pages.
`gh-pages` branch contains the website source, including configuration, website layout, and styling.
Markdown documents are automatically populated to `gh-pages` from the main branch and should be authored there.
Changes to the Aristi.io website layout and styling should be checked out from the `gh-pages` branch and  PRs should be created against `gh-pages`.

### Local website authoring and testing


Thanks for contributing!