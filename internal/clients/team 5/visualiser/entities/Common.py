"""
Common functions between entities.
"""

class Drawable:
    def __init__(self, x, y) -> None:
        self.x = x
        self.y = y

    def update_position(self, x:int, y:int) -> None:
        """
        Update the position of the agent.
        """
        self.x = x
        self.y = y
