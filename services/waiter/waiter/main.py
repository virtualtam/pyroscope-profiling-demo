from fastapi import FastAPI, Request, Response
from waiter.config import (
    LISTEN_PORT,
    LISTEN_ADDR,
    PYROSCOPE_ADDR,
)
import time
from waiter.logging import configure_logging
from asgi_correlation_id import CorrelationIdMiddleware
from asgi_correlation_id.context import correlation_id
import pyroscope
import uvicorn
import structlog
from uvicorn.protocols.utils import get_path_with_query_string
from waiter.api.router import api_router

APP_NAME = "demo.waiter"

logger = structlog.getLogger(__name__)

configure_logging()

access_logger = structlog.stdlib.get_logger("api.access")

app = FastAPI()

api = FastAPI(
    title="Waiter",
    description="Welcome to Waiter's API documentation!",
    root_path="/api/v1",
)


@app.middleware("http")
async def log_requests(request: Request, call_next):
    structlog.contextvars.clear_contextvars()
    # These context vars will be added to all log entries emitted during the request
    request_id = correlation_id.get()
    structlog.contextvars.bind_contextvars(request_id=request_id)

    start_time = time.perf_counter_ns()
    # If the call_next raises an error, we still want to return our own 500 response,
    # so we can add headers to it (process time, request ID...)
    response = Response(status_code=500)
    try:
        response = await call_next(request)
    except Exception:
        # TODO: Validate that we don't swallow exceptions (unit test?)
        structlog.stdlib.get_logger("api.error").exception("Uncaught exception")
        raise
    finally:
        process_time = time.perf_counter_ns() - start_time
        status_code = response.status_code
        url = get_path_with_query_string(request.scope)
        client_host = request.client.host
        client_port = request.client.port
        http_method = request.method
        http_version = request.scope["http_version"]
        # Recreate the Uvicorn access log format, but add all parameters as structured information
        access_logger.info(
            f"""{client_host}:{client_port} - "{http_method} {url} HTTP/{http_version}" {status_code}""",
            http={
                "url": str(request.url),
                "status_code": status_code,
                "method": http_method,
                "request_id": request_id,
                "version": http_version,
            },
            network={"client": {"ip": client_host, "port": client_port}},
            duration=process_time,
        )
        response.headers["X-Process-Time"] = str(process_time / 10**9)
        return response


app.add_middleware(CorrelationIdMiddleware)


@app.get("/")
async def index():
    return {"status": "ok"}


# we add all API routes to the Web API framework
api.include_router(api_router)

app.mount("/api/v1", app=api)

if __name__ == "__main__":
    if PYROSCOPE_ADDR != "":
        pyroscope.configure(
            app_name=APP_NAME,
            server_address=PYROSCOPE_ADDR,
        )

    uvicorn.run(app, host=LISTEN_ADDR, port=LISTEN_PORT, log_config=None)
