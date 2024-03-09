from fastapi import APIRouter

api_router_v1 = APIRouter()


@api_router_v1.get("/menus")
async def get_menus():
    return {"message": "menus v1"}
