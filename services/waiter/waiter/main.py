from fastapi import FastAPI
import pyroscope
import uvicorn

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
        app_name="demo.waiter",
        server_address="http://127.0.0.1:4040",  # replace this with the address of your pyroscope server
    )

    uvicorn.run(api, host="0.0.0.0", port=8081)
