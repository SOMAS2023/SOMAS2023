"""
Helper functions for the visualiser
"""
import random

def pick_random_colour(colourDict) -> str:
    """
    Pick a random colour from the list of colours
    """
    return random.choice(list(colourDict.values()))
