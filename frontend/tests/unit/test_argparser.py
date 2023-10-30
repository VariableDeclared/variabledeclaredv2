import sys
import os
import unittest

TEST_DIR = os.path.dirname(os.path.abspath(__name__))
sys.path.append(f"{TEST_DIR}/../../../")

from frontend.argparser import (
    VDotComArgParser,
    VDotComArg
)

class TestArgparser(unittest.TestCase):
    def test_argparser_singleton(self):
        inst1 = VDotComArgParser()
        inst2 = VDotComArgParser()

        assert id(inst1) == id(inst2)

class TestVDotComArgs(unittest.TestCase):
    def test_arg(self):
        arg = VDotComArg("test", "myaction", "BOO")

        assert arg.arg == "test"

    def test_arg_type(self):
        arg = VDotComArg("test", "myaction", "BOO")

        assert type(arg.arg) == str

    def test_action(self):
        arg = VDotComArg("test", "myaction", "BOO")

        assert arg.action == "myaction"

    def test_action_type(self):
        arg = VDotComArg("test", "myaction", "BOO")

        assert type(arg.action) == str
    def test_default(self):
        arg = VDotComArg("test", "myaction", "BOO")

        assert arg.default == "BOO"

    def test_default_type(self):
        arg = VDotComArg("test", "myaction", "BOO")

        assert type(arg.default) == str

if __name__ == '__main__':
    unittest.main()