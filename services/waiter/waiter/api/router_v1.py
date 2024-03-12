import requests
import logging
from fastapi import APIRouter, HTTPException
from http import HTTPStatus
from pydantic import ValidationError
from waiter.config import COOK_API_BASE_URL
from waiter.api.models import Menu

api_router_v1 = APIRouter()

logger = logging.getLogger(__name__)


@api_router_v1.get("/restaurant/{restaurant_id}/menu", response_model=Menu)
async def get_menu(restaurant_id: int):
    url = COOK_API_BASE_URL + "/v1/restaurant/" + str(restaurant_id) + "/menu"

    logger.info(f"requesting {url}")

    try:
        r = requests.get(url)
    except Exception as err:
        logger.error(err)

    if r.status_code != HTTPStatus.OK:
        raise HTTPException(status_code=r.status_code)

    return r.json()
