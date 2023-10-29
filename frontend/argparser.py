import argparse
import types

class VDotComArg:
    def __init__(self, arg: str, action: str = None, default: str = None, **kwargs) -> None:
        self._arg = arg
        self._action = action
        self._default = default
        # argparse kwargs.
        self._kwargs=kwargs
    @property
    def arg(self) -> str:
        return self._arg
    @property
    def action(self) -> str:
        return self._action
    @property
    def default(self) -> str or None:
        return self._default
    @property
    def kwargs(self) -> dict[str, type]:
        return self._kwargs

class VDotComArgParser:
    _args = [
        VDotComArg("--debug", "store_true"),
        VDotComArg("--port", default="8080")
    ]
    def __new__(cls):
        if not hasattr(cls, 'instance'):
            cls.instance = super(VDotComArgParser, cls).__new__(cls)
        return cls.instance

    def __init__(self, args: list[VDotComArg] = None) -> None:
        self._parser = argparse.ArgumentParser("VDotCom Frontend.")
        if args is not None:
            self._args = args
        self.setup_args()
    
    def setup_args(self) -> None:
        for arg in self._args:
            self._parser.add_argument(arg.arg, action=arg.action)
    
    @property
    def parser(self) -> argparse.ArgumentParser:
        return self._parser


