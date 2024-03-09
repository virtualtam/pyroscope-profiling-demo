import structlog
import time
from asgi_correlation_id.context import correlation_id
from fastapi import Request, Response
from uvicorn.protocols.utils import get_path_with_query_string


async def log_requests(request: Request, call_next):
    access_logger = structlog.stdlib.get_logger("api.access")

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
