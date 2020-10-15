class InvalidUsage(Exception):
    """Creates an exception to be returned by the api"""
    status_code = 400

    def __init__(self, message, status_code=None, meta=None, payload=None):
        Exception.__init__(self)
        self.message = message
        if status_code is not None:
            self.status_code = status_code
        self.meta = meta
        self.payload = payload

    def to_dict(self):
        """Creates a dict from the error message to be returned as json"""
        error_message = dict(self.payload or ())
        error_message['error'] = self.message
        if self.meta and self.meta is not None:
            error_message['meta'] = self.meta
        return error_message
