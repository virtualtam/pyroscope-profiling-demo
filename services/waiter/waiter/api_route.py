from fastapi import APIRouter
from waiter.api.router_v1 import api_router_v1
from waiter.api.router_v2 import api_router_v2


api_router = APIRouter()

api_router.include_router(
    api_router_v1,
    prefix="/v1",
)

api_router.include_router(
    api_router_v2,
    prefix="/v2",
)
