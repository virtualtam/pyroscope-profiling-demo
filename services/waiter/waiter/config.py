import logging
from starlette.config import Config

config = Config(env_prefix="WAITER_")

LOG_LEVEL = config("LOG_LEVEL", default=logging.WARNING)
LISTEN_PORT = config("LISTEN_PORT", default=8081, cast=int)
LISTEN_ADDR = config("LISTEN_ADDR", default="0.0.0.0")
PYROSCOPE_ADDR = config("PYROSCOPE_ADDR", default="")
COOK_API_BASE_URL = config("COOK_ADDR", default="http://cook:8080/api")
