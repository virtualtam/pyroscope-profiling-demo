FROM python:3.11-slim


ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    POETRY_NO_INTERACTION=1 \
    POETRY_VIRTUALENVS_IN_PROJECT=true

ARG POETRY_VERSION=1.8.2
RUN pip3 install "poetry==${POETRY_VERSION}"

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

USER waiter
WORKDIR /opt/waiter

COPY pyproject.toml poetry.lock README.md ./

RUN poetry install

COPY . .

EXPOSE 8081

ENTRYPOINT ["poetry", "run", "python", "-m", "waiter.main"]
