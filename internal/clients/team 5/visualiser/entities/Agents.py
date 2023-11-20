"""
Logic for handling the agents in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
import math
import pygame
import pygame_gui
from visualiser.util.Constants import AGENT
from visualiser.entities.Common import Drawable

class Agent(Drawable):
    def __init__(self, x:int, y:int, colour:pygame.color, groupID, agentID) -> None:
        super().__init__(x, y)
        self.colour = colour
        self.radius = AGENT["SIZE"]
        self.groupID = groupID
        self.id = agentID

    def check_collision(self, mouseX:int, mouseY:int, offsetX:int, offsetY:int, zoom:float) -> bool:
        """
        Check if the mouse click intersects with the agent.
        """
        # Check if the mouse click is within the agent's radius
        zoomedRadius = int(self.radius * zoom)
        distance = math.sqrt((self.trueX - mouseX) ** 2 + (self.trueY - mouseY) ** 2)
        # Check if the distance is within the agent's radius
        return distance <= zoomedRadius

    def draw(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Draw the agent on the screen.
        """
        zoomedRadius = int(self.radius * zoom)
        self.trueX = int(self.x * zoom + offsetX)
        self.trueY = int(self.y * zoom + offsetY)
        pygame.draw.circle(screen, AGENT["LINE_COLOUR"], (self.trueX, self.trueY), zoomedRadius+max(AGENT["LINE_WIDTH"]*zoom, 1))
        pygame.draw.circle(screen, self.colour, (self.trueX, self.trueY), zoomedRadius)
        # Draw group ID
        font = pygame.font.SysFont("Arial", int(AGENT["FONT_SIZE"]*zoom))
        if self.colour == "White":
            text = font.render(str(self.groupID), True, (0, 0, 0))
        else:
            text = font.render(str(self.groupID), True, (255, 255, 255))
        textRect = text.get_rect()
        textRect.center = (self.trueX, self.trueY)
        screen.blit(text, textRect)
        self.overlay = self.update_overlay(zoom)
