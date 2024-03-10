# Importing the Food.com dataset
This demo uses data from the [Food.com Recipes and Interactions](https://www.kaggle.com/datasets/shuyangli94/food-com-recipes-and-user-interactions)
dataset to populate the PostgreSQL database.

We will use the following data:

- a list of recipes from [Food.](https://www.food.com/);
- a list of the necessary ingredients for each recipe.

## Get the dataset

1. Login to [Kaggle](https://www.kaggle.com/)
2. [Download](https://www.kaggle.com/datasets/shuyangli94/food-com-recipes-and-user-interactions?resource=download) the full dataset into the `dataset` directory
3. Extract the `archive.zip` archive

## Install script dependencies

Create and activate a Python3 virtualenv:

```shell
$ python3 -m venv .venv
$ source .venv/bin/activate
```

Install dependencies:

```shell
$ pip install -r requirements.txt
```

## Prepare data

The `prepare_data.py` script will process the dataset to output a JSON file with curated data,
that can then be imported into the PostgreSQL database.

Run the script with:

```shell
$ ./prepare_data.py RAW_recipes.csv recipes.out.json
```
