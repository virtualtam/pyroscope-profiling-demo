from pydantic import BaseModel
from typing import List


class Ingredient(BaseModel):
    id: int
    name: str


class Dish(BaseModel):
    id: int
    name: str
    price: float
    ingredients: List[Ingredient]


class Menu(BaseModel):
    id: int
    dishes: List[Dish]
