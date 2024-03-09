#!/usr/bin/env python3
from argparse import ArgumentParser
from ast import literal_eval
import logging

import pandas as pd


def main():
    parser = ArgumentParser()
    parser.add_argument(
        "input",
        help="CSV input file for raw recipes data",
    )
    parser.add_argument(
        "output",
        help="JSON output file for transformed recipes data",
    )

    args = parser.parse_args()

    recipes = prepare_data(args.input)
    recipes.to_json(args.output, orient="records")


def prepare_data(recipe_csv: str) -> pd.DataFrame:
    ingredient_map = pd.read_pickle("ingr_map.pkl")
    recipes = pd.read_csv(recipe_csv)

    recipes.ingredients = recipes.ingredients.apply(literal_eval)

    recipes.ingredients = recipes.ingredients.map(
        lambda ingredients: replace_ingredients(ingredient_map, ingredients)
    )

    return recipes[["name", "ingredients"]]


def replace_ingredients(
    ingredient_map: pd.DataFrame, ingredients: list[str]
) -> list[str]:
    return [
        replace_ingredient(ingredient_map, ingredient) for ingredient in ingredients
    ]


def replace_ingredient(ingredient_map: pd.DataFrame, ingredient: str) -> str:
    try:
        return (
            ingredient_map["replaced"]
            .loc[ingredient_map["raw_ingr"] == ingredient]
            .values[0]
        )
    except IndexError:
        logging.debug(f"no replacement for: {ingredient}")
        return ingredient


if __name__ == "__main__":
    main()
