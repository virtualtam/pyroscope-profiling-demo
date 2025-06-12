FROM python:3.13-slim
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    UV_COMPILE_BYTECODE=1 \
    UV_LINK_MODE=copy \
    WAITER_LOG_LEVEL="INFO" \
    WAITER_LISTEN_PORT=8081 \
    WAITER_LISTEN_ADDR="0.0.0.0" \
    WAITER_PYROSCOPE_ADDR="http://pyroscope:4040"

RUN groupadd \
    --gid 1000 \
    waiter \
    && useradd \
    --create-home \
    --home-dir /var/lib/waiter \
    --shell /bin/bash \
    --uid 1000 \
    --gid waiter \
    waiter

RUN mkdir -p /opt/waiter \
    && chown -R waiter:waiter /opt/waiter

WORKDIR /opt/waiter

RUN --mount=type=cache,target=/root/.cache/uv \
    --mount=type=bind,source=uv.lock,target=uv.lock,readonly \
    --mount=type=bind,source=pyproject.toml,target=pyproject.toml,readonly \
    uv sync --locked --no-install-project --no-dev

COPY . .
RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --locked --no-dev

USER waiter

EXPOSE 8081

ENTRYPOINT ["uv", "run", "python", "-m", "waiter.main"]
