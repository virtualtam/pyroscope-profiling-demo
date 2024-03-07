from starlette.config import Config

config = Config(env_prefix="WAITER_")

LISTEN_PORT = config("LISTEN_PORT", default=8081, cast=int)
LISTEN_ADDR = config("LISTEN_ADDR", default="0.0.0.0")
PYROSCOPE_ADDR = config("PYROSCOPE_ADDR")
