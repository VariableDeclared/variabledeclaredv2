

class BlogApp:
    def __init__(self, title: str, subheading: str) -> None:
        self._app_title = title
        self._app_subheading = subheading

    @property
    def title(self):
        return self._app_title
    @property
    def subheading(self):
        return self._app_subheading