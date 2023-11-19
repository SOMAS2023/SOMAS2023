"""
Logic for handling the agents in the visualiser
"""
import pygame
import pygame_gui
from visualiser.util.Constants import DIM

class Agent:
    def __init__(self, x:int, y:int, colour:pygame.color) -> None:
        self.x = x
        self.y = y
        self.colour = colour
        self.radius = DIM["AGENT_SIZE"]

    def draw(self,  screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float):
        """
        Draw the agent on the screen.
        """
        zoomedRadius = int(self.radius * zoom)
        zoomedX = int(self.x * zoom + offsetX)
        zoomedY = int(self.y * zoom + offsetY)
        pygame.draw.circle(screen, self.colour, (zoomedX, zoomedY), zoomedRadius)

    def update_position(self, x:int, y:int) -> None:
        """
        Update the position of the agent.
        """
        self.x = x
        self.y = y
