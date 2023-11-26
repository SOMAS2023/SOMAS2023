"""
Helper functions for the visualiser
"""
import random

def pick_random_colour(colourDict:dict) -> str:
    """
    Pick a random colour from the list of colours
    """
    return random.choice(list(colourDict.values()))

def make_center(elementSize:tuple, containerSize:tuple) -> tuple:
    """
    Center an element in a container
    """
    elementWidth, elementHeight = elementSize
    containerWidth, containerHeight = containerSize
    x = (containerWidth - elementWidth) // 2
    y = (containerHeight - elementHeight) // 2
    return (x, y)
