# Blinders

Monorepo, microservice back-end project named Blinders for Peakee - a language learning platform.

This project is a monorepo microservice back-end written in Golang and Python. It's using Terraform for setup provisioning, all the services are hosted on AWS, the main AWS service used is AWS lambda.

## Setup

### Python

-   [x] Monorepo structure
-   [x] Setup linter and formatter
-   [x] Testing setup
-   [x] AWS lambda function
-   [x] Terraform deployment
-   [x] CI Github Actions
-   [ ] CD Github Actions

### Golang

-   [x] Monorepo structure
-   [x] Setup linter and formatter
-   [ ] Testing setup
-   [ ] AWS lambda function
-   [ ] Terraform deployment
-   [x] CI Github Actions
-   [ ] CD Github Actions

## Python development setup

### Tools

-   Code formatter: [Black](https://github.com/psf/black)
-   Code linter: [Flake8](https://flake8.pycqa.org/en/latest/user/index.html) [isort](https://github.com/PyCQA/isort), and [pylint](https://pypi.org/project/pylint/) for just checking public artifacts are documented or not
-   Type checking: [Pyright](https://github.com/microsoft/pyright#static-type-checker-for-python)

### Setup steps

Use python 3.10 as base version of python, recommend to use local python environment using [conda](https://www.anaconda.com/)

```shell
conda create --prefix ./.venv/ python==3.10 # Initialize repo virtual environment
conda activate ./.venv # Activate venv
```

We're using [poetry](https://python-poetry.org/) package manager because of rich dependencies management features

```shell
pip install poetry
```

Install packages

```shell
poetry install
```

If not using `poetry`

```shell
pip install -e .
```

## Go development setup

-   Using [golangci-lint](https://golangci-lint.run/) to manage all linter, formatter and setup ci, detail configs in `golangci.yml`

## References

-   Awesome monorepo [github.com/korfuri/awesome-monorepo](https://github.com/korfuri/awesome-monorepo)
-   Setup Python monorepo
    -   [tweag.io/blog/2023-04-04-python-monorepo-1](https://www.tweag.io/blog/2023-04-04-python-monorepo-1/)
    -   [tweag.io/blog/2023-07-13-python-monorepo-2](https://www.tweag.io/blog/2023-07-13-python-monorepo-2/)
    -   [medium.com/opendoor-labs/our-python-monorepo-d34028f2b6fa](medium.com/opendoor-labs/our-python-monorepo-d34028f2b6fa)
-   Use pyproject.toml [peps.python.org/pep-0518](https://peps.python.org/pep-0518/)
-   Poetry - python package manager [python-poetry.org](https://python-poetry.org/docs/)
-   Python namespace [packaging.python.org/en/latest/guides/packaging-namespace-packages](https://packaging.python.org/en/latest/guides/packaging-namespace-packages/)
-   Python editable mode for development [pip.pypa.io/en/stable/topics/local-project-installs](https://pip.pypa.io/en/stable/topics/local-project-installs/)
-   Building lambda with Poetry [chariotsolutions.com/blog/post/building-lambdas-with-poetry/](https://chariotsolutions.com/blog/post/building-lambdas-with-poetry/)
-   Effective Go [go.dev/doc/effective_go](https://go.dev/doc/effective_go)
-   Go setup formatter and linter [medium.com/cp-massive-programming/golang-automate-formatting-and-linting-via-pre-commit-c43740065c2e](https://medium.com/cp-massive-programming/golang-automate-formatting-and-linting-via-pre-commit-c43740065c2e)
