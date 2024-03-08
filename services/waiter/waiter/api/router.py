from fastapi import APIRouter

api_router = APIRouter()


@api_router.get("/menus")
async def get_menus():
    return {}
