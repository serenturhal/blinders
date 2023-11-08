isort --check-only $(git ls-files "*.py")
black --check $(git ls-files "*.py")
pylint $(git ls-files "*.py")
flake8 $(git ls-files "*.py")
