"""
Logic for handling bikes in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
import pygame
import pygame_gui
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
        Draw the bike on the screen.
        """
        noAgents = len(self.agentList.values())
        rectWidth = BIKE["SIZE"] * zoom
        # Adjust the height calculation to provide enough space for all agents
        rectHeight = BIKE["SIZE"] * zoom * max(noAgents, 1) * 2.5
        rectX = int(self.x * zoom + offsetX - rectWidth / 2)
        # Adjust the Y position to start from the top
        rectY = int(self.y * zoom + offsetY - rectHeight)
        lineWidth = int(BIKE["LINE_WIDTH"] * zoom)
        # Draw the filled rectangle
        pygame.draw.rect(screen, BIKE["COLOUR"], (rectX, rectY, rectWidth, rectHeight))
        # Draw the outline rectangle
        pygame.draw.rect(screen, BIKE["LINE_COLOUR"], (rectX, rectY, rectWidth, rectHeight), max(1, lineWidth))
        # Draw the agents, starting from the top of the bike
        agentSpacing = rectHeight / max(noAgents, 1)
        agentSize = AGENT["SIZE"] * zoom  # Assuming AGENT["SIZE"] is defined
        for index, agent in enumerate(self.agentList.values()):
            # Calculate the Y position for each agent
            agentY = rectY + 2*agentSize - (rectHeight / 4) + (index * agentSpacing)
            agent.draw(screen, offsetX, agentY, zoom)

    def set_agents(self, agentJson:dict) -> None:
        """
        Set the agents that are in the bike
        """
        self.agentList = dict()
        for agent, data in agentJson.items():
            self.agentList[agent] = Agent(self.x, self.y, data["colour"], data["groupID"], agent)
