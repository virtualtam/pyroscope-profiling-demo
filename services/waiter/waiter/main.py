import pyroscope
import structlog
import uvicorn
from asgi_correlation_id import CorrelationIdMiddleware
from fastapi import FastAPI
from starlette.middleware.base import BaseHTTPMiddleware
from waiter.api_route import api_router
from waiter.config import (
    LISTEN_PORT,
    LISTEN_ADDR,
    PYROSCOPE_ADDR,
)
from waiter.logging import configure_logging
from waiter.middleware import log_requests


APP_NAME = "demo.waiter"

logger = structlog.getLogger(__name__)

configure_logging()


app = FastAPI(
    docs_url=None,
)

app.add_middleware(BaseHTTPMiddleware, dispatch=log_requests)
# Need to add at the end as middlewares are applied in reverse order
app.add_middleware(CorrelationIdMiddleware)


@app.get("/")
async def index():
    return {"status": "ok"}


api = FastAPI(
    title="Waiter",
    description="Welcome to Waiter's API documentation!",
    root_path="/api",
)

# we add all API routes to the Web API framework
api.include_router(api_router)

app.mount("/api", app=api)

if __name__ == "__main__":
    if PYROSCOPE_ADDR != "":
        pyroscope.configure(
            app_name=APP_NAME,
            server_address=PYROSCOPE_ADDR,
        )

    uvicorn.run(app, host=LISTEN_ADDR, port=LISTEN_PORT, log_config=None)
