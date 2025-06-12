#!/usr/bin/env -S uv run --script
# /// script
# dependencies = [
#   "pandas==2.3.0",
#   "pyarrow==20.0.0",  # For faster CSV reading and better memory usage
# ]
# ///
from argparse import ArgumentParser
from ast import literal_eval
import logging
import multiprocessing
from typing import Dict, List

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
    parser.add_argument(
        "--chunk-size",
        type=int,
        default=10000,
        help="Size of chunks for processing CSV data",
    )
    parser.add_argument(
        "--limit",
        type=int,
        default=None,
        help="Limit number of recipes to process",
    )
    parser.add_argument(
        "--workers",
        type=int,
        default=multiprocessing.cpu_count(),
        help="Number of workers for parallel processing",
    )

    args = parser.parse_args()

    recipes = prepare_data(args.input, args.limit, args.chunk_size, args.workers)
    recipes.to_json(args.output, orient="records")


def prepare_data(
    recipe_csv: str, limit: int, chunk_size: int, workers: int
) -> pd.DataFrame:
    ingredient_df = pd.read_pickle("ingr_map.pkl")
    ingredient_map = dict(zip(ingredient_df["raw_ingr"], ingredient_df["replaced"]))

    chunks = pd.read_csv(recipe_csv, chunksize=chunk_size, nrows=limit)

    with multiprocessing.Pool(processes=workers) as pool:
        processed_chunks = pool.starmap(
            process_chunk, [(chunk, ingredient_map) for chunk in chunks]
        )

    return pd.concat(processed_chunks, ignore_index=True)


def process_chunk(chunk: pd.DataFrame, ingredient_map: Dict[str, str]) -> pd.DataFrame:
    chunk.ingredients = chunk.ingredients.apply(literal_eval)
    chunk.ingredients = chunk.ingredients.map(
        lambda ingredients: replace_ingredients(ingredient_map, ingredients)
    )
    return chunk[["name", "ingredients"]]


def replace_ingredients(
    ingredient_map: Dict[str, str], ingredients: List[str]
) -> List[str]:
    return [
        replace_ingredient(ingredient_map, ingredient) for ingredient in ingredients
    ]


def replace_ingredient(ingredient_map: Dict[str, str], ingredient: str) -> str:
    if ingredient in ingredient_map:
        return ingredient_map[ingredient]

    logging.debug(f"no replacement for: {ingredient}")
    return ingredient


if __name__ == "__main__":
    main()
