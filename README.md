<h1 align="center">
  <br>
  <a href="http://github.com/ghostsquad/s3-file-explorer"><img src="./docs/assets/cloud-service.png" alt="github.com/ghostsquad/s3-file-explorer" width="200px" /></a>
  <br>
  S3 File Explorer
  <br>
</h1>

<p align="center">
  <a href="#introduction">Introduction</a> •
  <a href="#getting-started">Getting Started</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#roadmap">Roadmap</a>
</p>

## Introduction

A Simple S3 File Explorer. Make a request to `/p/` to see what's in the configured bucket.

## Getting Started

```shell
task run
```

In a separate shell

```shell
task http:paths

task http:metrics
```

### Configuration

Configuration is done via environment variables. Standard AWS SDK Environment variables supported, as well as OIDC/EC2 authentication methods.

| Variable       | Required | Default                              | Description                                                    |
|----------------|----------|--------------------------------------|----------------------------------------------------------------|
| `AWS_BUCKET`   | no       | `github-ghostsquad-s3-file-explorer` | The AWS Bucket to explorer                                     |
| `PORT`         | no       | `8080`                               | The listen port                                                |
| `BIND_ADDRESS` | no       |                                      | Configured to listen on 127.0.0.1, this may not work in Docker |

## Contributing

```shell
brew install asdf
asdf plugin-add task https://github.com/particledecay/asdf-task.git
asdf plugin add python

asdf install

task test
```

## Roadmap

- [ ] Make a roadmap

## Attribution

<a href="https://www.flaticon.com/free-icons/function" title="function icons">Function icons created by Freepik - Flaticon</a>