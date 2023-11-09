from blinders.dictionary_core import handle_dictionary


def call_handle_dictionary():
    return "called " + handle_dictionary()


def lambda_handler(event, context):
    print("dictionary on event", event, context)
    return call_handle_dictionary()
