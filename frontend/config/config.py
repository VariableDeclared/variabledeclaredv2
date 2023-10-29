import json
import os
import pathlib
from base64 import b64decode

class VDotComConfig:
    _epilogue: str

    def __init__(self, epilogue: str) -> None:
        self._epilogue = epilogue
    
    @staticmethod
    def load_json(loaded_json: dict[type,type]):
        # TODO: Scales terribly. Need to revamp whole deserialisation 
        # logic.
        return VDotComConfig(loaded_json.get('epilogue', None))
    
    @property
    def epilogue(self):
        return b64decode(self._epilogue).decode('utf-8')



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