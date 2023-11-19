"""
Logic for handling the agents in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
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
        self.agentID = agentID

    def draw(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Draw the agent on the screen.
        """
        zoomedRadius = int(self.radius * zoom)
        zoomedX = int(self.x * zoom + offsetX)
        zoomedY = int(self.y * zoom + offsetY)
        pygame.draw.circle(screen, AGENT["LINE_COLOUR"], (zoomedX, zoomedY), zoomedRadius+max(AGENT["LINE_WIDTH"]*zoom, 1))
        pygame.draw.circle(screen, self.colour, (zoomedX, zoomedY), zoomedRadius)
        # Draw group ID
        font = pygame.font.SysFont("Arial", int(AGENT["FONT_SIZE"]*zoom))
        if self.colour == "White":
            text = font.render(str(self.groupID), True, (0, 0, 0))
        else:
            text = font.render(str(self.groupID), True, (255, 255, 255))
        textRect = text.get_rect()
        textRect.center = (zoomedX, zoomedY)
        screen.blit(text, textRect)
