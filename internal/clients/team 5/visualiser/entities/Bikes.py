"""
Logic for handling bikes in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
import math
import pygame
from visualiser.util.Constants import AGENT, BIKE
from visualiser.entities.Agents import Agent
from visualiser.entities.Common import Drawable

class Bike(Drawable):
    def __init__(self, x:int, y:int, bikeID) -> None:
        super().__init__(x, y)
        self.bikeID = bikeID
        self.agentList = dict()

    def draw(self, screen, offsetX, offsetY, zoom) -> None:
        """
        Draw the bike and agents onto the screen.
        """
        noAgents = len(self.agentList.values())
        # Determine the grid size
        gridSize = min(3, max(1, int(math.ceil(math.sqrt(noAgents)))))
        agentPadding = AGENT["PADDING"] * zoom  # Padding around agents
        agentSize = AGENT["SIZE"] * 2 * zoom
        squareSide = gridSize * agentSize + ((gridSize+1) * agentPadding)
        # Calculate bike's position and size
        rectX = int(self.x * zoom + offsetX)
        rectY = int(self.y * zoom + offsetY)
        # Draw the bike square
        pygame.draw.rect(screen, BIKE["COLOUR"], (rectX, rectY, squareSide, squareSide))
        # Draw the agents within the bike
        for index, agent in enumerate(self.agentList.values()):
            # Calculate agent's position within the grid
            row = index // gridSize
            col = index % gridSize
            agentX = rectX - 3 * agentSize + ((agentSize + agentPadding) * (col + 1))
            agentY = rectY - 3 * agentSize + ((agentSize + agentPadding) * (row + 1))
            agent.draw(screen, agentX, agentY, zoom)

    def set_agents(self, agentJson:dict) -> None:
        """
        Set the agents that are in the bike
        """
        self.agentList = dict()
        for agent, data in agentJson.items():
            self.agentList[agent] = Agent(self.x, self.y, data["colour"], data["groupID"], agent)
