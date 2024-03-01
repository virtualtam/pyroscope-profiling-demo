from fastapi import FastAPI
import pyroscope
import requests

app = FastAPI()

pyroscope.configure(
    application_name="demo.waiter",
    server_address="http://localhost:4040",
)


@app.get("/")
async def root():
    return {"message": "Hello World"}


@app.get("/menu")
async def root():
    resp = requests.get("http://localhost:8080")
    return resp.json()
