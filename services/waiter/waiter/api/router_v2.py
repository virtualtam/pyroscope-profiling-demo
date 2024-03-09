from fastapi import APIRouter

api_router_v2 = APIRouter()


@api_router_v2.get("/menus")
async def get_menus():
    return {"message": "menus v2"}
