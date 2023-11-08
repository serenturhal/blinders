import os


def handle_dictionary():
    print("called handle dictionary", os.getcwd())
    return "handle dictionary"


__all__ = ["handle_dictionary"]