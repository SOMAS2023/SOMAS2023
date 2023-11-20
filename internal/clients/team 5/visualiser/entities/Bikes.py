"""
Logic for handling bikes in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
import math
import pygame
import pygame_gui
from visualiser.util.Constants import AGENT, BIKE
from visualiser.entities.Agents import Agent
from visualiser.entities.Common import Drawable

class Bike(Drawable):
    def __init__(self, x:int, y:int, bikeID) -> None:
        super().__init__(x, y)
        self.id = bikeID
        self.agentList = dict()
        self.squareSide = 0

    def draw(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Draw the bike and agents onto the screen.
        """
        noAgents = len(self.agentList.values())
        # Determine the grid size
        gridSize = min(3, max(1, int(math.ceil(math.sqrt(noAgents)))))
        agentPadding = AGENT["PADDING"] * zoom  # Padding around agents
        agentSize = AGENT["SIZE"] * 2 * zoom
        self.squareSide = gridSize * agentSize + ((gridSize+1) * agentPadding)
        # Calculate bike's position and size
        self.trueX = int(self.x * zoom + offsetX)
        self.trueY = int(self.y * zoom + offsetY)
        # Draw the bike square
        pygame.draw.rect(screen, BIKE["COLOUR"], (self.trueX, self.trueY, self.squareSide, self.squareSide))
        # Draw the agents within the bike
        for index, agent in enumerate(self.agentList.values()):
            # Calculate agent's position within the grid
            row = index // gridSize
            col = index % gridSize
            agentX = self.trueX - 3 * agentSize + ((agentSize + agentPadding) * (col + 1))
            agentY = self.trueY - 3 * agentSize + ((agentSize + agentPadding) * (row + 1))
            agent.draw(screen, agentX, agentY, zoom)
        for index, agent in enumerate(self.agentList.values()):
            agent.draw_overlay(screen)
        self.overlay = self.update_overlay(zoom)
        self.draw_overlay(screen)

    def set_agents(self, agentJson:dict) -> None:
        """
        Set the agents that are in the bike
        """
        self.agentList = dict()
        for agent, data in agentJson.items():
            self.agentList[agent] = Agent(self.x, self.y, data["colour"], data["groupID"], agent)

    def propagate_click(self, mouseX:int, mouseY:int, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Propagate the click to the agents within the bike
        """
        intersected = False
        for agent in self.agentList.values():
            if agent.click(mouseX, mouseY, offsetX, offsetY, zoom):
                intersected = True
        # If an agent was not interacted with, check bike
        if not intersected:
            self.click(mouseX, mouseY, offsetX, offsetY, zoom)

    def check_collision(self, mouseX: int, mouseY: int, offsetX: int, offsetY: int, zoom: float) -> bool:
        """
        Check if the mouse click intersects with the bike.
        """
        return (self.trueX <= mouseX <= self.trueX + self.squareSide) and \
               (self.trueY <= mouseY <= self.trueY + self.squareSide)
