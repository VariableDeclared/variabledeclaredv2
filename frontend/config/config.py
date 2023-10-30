import json
import os
import pathlib
from base64 import b64decode
from blog_app.blog import BlogApp

class VDotComConfig:
    _epilogue: str
    _app = None
    def __init__(self, epilogue: str, app: dict = None) -> None:
        self._epilogue = epilogue
        if app:
            self._app = BlogApp(
                app.get('title', None),
                app.get('subheading', None)
            )
        else:
            self._app = BlogApp(
                "My blog!",
                "Welcome!"
            )
    
    @staticmethod
    def load_json(loaded_json: dict[type,type]):
        # TODO: Scales terribly. Need to revamp whole deserialisation 
        # logic.
        app = loaded_json.get('app', None)
        config = VDotComConfig(loaded_json.get('epilogue', ''))
        if app:
            config = VDotComConfig(loaded_json.get('epilogue', ''), app=app)
        return config

    
    @property
    def epilogue(self) -> str:
        return b64decode(self._epilogue).decode('utf-8')
    @property
    def app(self) -> BlogApp:
        return self._app



class VDotComConfigLoader:
    _config_file: str = "vdotcom.json"
    _config: VDotComConfig = None
    _instances = {}
    def __call__(cls, *args, **kwargs):
        if cls not in cls._instances:
            cls._instances[cls] = super(VDotComConfigLoader, cls).__call__(*args, **kwargs)
        return cls._instances[cls]

    def __init__(self, config_file: str = None) -> None:
        if config_file:
            self._config_file = config_file
        self.load_config()

    def load_config(self):
        with open(self._config_file, 'r') as fh:
            loaded_json = json.loads(fh.read())
        self._config = VDotComConfig.load_json(loaded_json)

    @property
    def config(self) -> VDotComConfig:
        return self._config