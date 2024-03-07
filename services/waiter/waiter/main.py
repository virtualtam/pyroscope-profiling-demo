from fastapi import FastAPI
from .config import (
    LISTEN_PORT,
    LISTEN_ADDR,
    PYROSCOPE_ADDR,
)
import pyroscope
import uvicorn

APP_NAME = "demo.waiter"

api = FastAPI(
    title="Waiter",
    description="Welcome to Waiter's API documentation!",
    root_path="/api/v1",
    docs_url=None,
    openapi_url="/docs/openapi.json",
    redoc_url="/docs",
)

if __name__ == "__main__":
    pyroscope.configure(
        app_name=APP_NAME,
        server_address=PYROSCOPE_ADDR,
    )

    uvicorn.run(api, host=LISTEN_ADDR, port=LISTEN_PORT)
