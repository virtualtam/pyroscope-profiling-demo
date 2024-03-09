import requests
import logging
from fastapi import APIRouter, HTTPException
from http import HTTPStatus
from pydantic import ValidationError
from waiter.config import COOK_API_BASE_URL
from waiter.api.models import Menu
from functools import lru_cache

api_router_v2 = APIRouter()

logger = logging.getLogger(__name__)


@lru_cache(maxsize=16)
@api_router_v2.get("/restaurant/{restaurant_id}/menu")
async def get_menu(restaurant_id: int):
    url = COOK_API_BASE_URL + "/v1/restaurant/" + str(restaurant_id) + "/menu"

    logger.info(f"requesting {url}")

    try:
        r = requests.get(url)
    except Exception as err:
        logger.error(err)

    if r.status_code != HTTPStatus.OK:
        raise HTTPException(status_code=r.status_code)

    try:
        Menu.model_validate(r.json())
    except ValidationError as err:
        logger.error(err)

    return r.json()
