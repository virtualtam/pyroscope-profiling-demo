import requests
import logging
from fastapi import APIRouter, Depends, HTTPException
from http import HTTPStatus
from waiter.config import COOK_API_BASE_URL
from waiter.api.models import Menu
from functools import lru_cache

api_router_v2 = APIRouter()

logger = logging.getLogger(__name__)


@lru_cache(maxsize=16)
def get_menu_from_cook(restaurant_id: int):
    logger.info(get_menu_from_cook.cache_info())
    url = COOK_API_BASE_URL + "/v1/restaurant/" + str(restaurant_id) + "/menu"

    logger.info(f"requesting {url}")

    try:
        r = requests.get(url)
    except Exception as err:
        logger.error(err)

    if r.status_code != HTTPStatus.OK:
        raise HTTPException(status_code=r.status_code)

    return r.json()


@api_router_v2.get("/restaurant/{restaurant_id}/menu", response_model=Menu)
async def get_menu(restaurant_id: int, menu=Depends(get_menu_from_cook)):
    return menu
