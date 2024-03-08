from pydantic import BaseModel
from typing import List
import uuid


class Ingredient(BaseModel):
    id: uuid
    name: str


class Dish(BaseModel):
    id: uuid
    name: str
    ingredients: List[Ingredient]


class Menu(BaseModel):
    id: uuid
    name: str
    price: float
    dishes: List[Dish]
