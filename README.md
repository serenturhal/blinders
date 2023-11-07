# Monorepo APIs for Peakee

## Setup

-   [ ] Monorepo structure
-   [ ] Go lambda function
-   [ ] Python lambda function
-   [ ] Setup linter and formatter for Go
-   [ ] Setup linter and formatter for Python
-   [ ] Testing setup for Go
-   [ ] Testing setup for Python
-   [ ] CI Github Actions
-   [ ] CD Github Actions

## Python development setup

### Description

Code formatter: [Black](https://github.com/psf/black)
Code linter: [Flake8](https://flake8.pycqa.org/en/latest/user/index.html) [isort](https://github.com/PyCQA/isort), and [pylint](https://pypi.org/project/pylint/) for just checking public artifacts are documented or not
Type checking: [Pyright](https://github.com/microsoft/pyright#static-type-checker-for-python)

### Setup steps

Use python 3.10 as base version of python, recommend to use local python environment using `conda`

```shell
conda create --prefix ./.venv/ python==3.10
```

Install pinned pip

```shell
pip install -r pip-requirements.txt
```

Install shared development dependencies

```shell
pip install -r dev-requirements.txt
```

## References

-   Setup Python monorepo [tweag.io/blog/2023-04-04-python-monorepo-1/](https://www.tweag.io/blog/2023-04-04-python-monorepo-1/) [tweag.io/blog/2023-07-13-python-monorepo-2/](https://www.tweag.io/blog/2023-07-13-python-monorepo-2/)
-   Use pyproject.toml [peps.python.org/pep-0518/](https://peps.python.org/pep-0518/)
-   Awesome monorepo [https://github.com/korfuri/awesome-monorepo](https://github.com/korfuri/awesome-monorepo)
