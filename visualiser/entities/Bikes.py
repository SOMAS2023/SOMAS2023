"""
Logic for handling bikes in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
import math
import pygame
import pygame_gui
from visualiser.util.Constants import AGENT, BIKE, PRECISION, COORDINATESCALE, GOVERNANCE
from visualiser.entities.Agents import Agent
from visualiser.entities.Common import Drawable

class Bike(Drawable):
    def __init__(self, bikeid:str, jsonData:dict, agentData:dict, nextPos:float, nextOrient:float, nextRound:dict) -> None:
        super().__init__(bikeid, jsonData)
        self.agentList = dict()
        self.update_bike(jsonData, agentData, nextPos, nextOrient, nextRound)

    def update_bike(self, jsonData:dict, agentData:dict, nextPos:float, nextOrient:float, nextRound:dict) -> None:
        """
        Update the bike's properties
        """
        super().update_entity(jsonData)
        self.agentList = dict()
        self.agentData = jsonData["agent_ids"]
        self.squareSide = 0
        self.orientation = nextOrient
        self.nextPos = nextPos
        self.governance, self.colour = GOVERNANCE[jsonData["governance"]]
        properties = {
            "Acceleration" : round(jsonData["physical_state"]["acceleration"], PRECISION),
            "Velocity" : round(jsonData["physical_state"]["velocity"], PRECISION),
            "Mass" : jsonData["physical_state"]["mass"],
            "Governance" : self.governance,
        }
        self.properties.update(properties)
        self.set_agents(agentData, nextRound)

    def draw(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Draw the bike and agents onto the screen.
        """
        noAgents = len(self.agentList.values())
        # Determine the grid size
        gridSize = max(1, int(math.ceil(math.sqrt(noAgents))))
        agentPadding = AGENT["PADDING"] * zoom  # Padding around agents
        agentSize = (AGENT["SIZE"] + AGENT["LINE_WIDTH"]) * 2 * zoom
        self.squareSide = gridSize * agentSize + ((gridSize+1) * agentPadding)
        # Calculate bike's position and size
        self.trueX = int(self.x * COORDINATESCALE * zoom + offsetX - self.squareSide/2)
        self.trueY = int(self.y * COORDINATESCALE * zoom + offsetY - self.squareSide/2)
        # Draw the bike square
        border = pygame.Surface(((2*BIKE["LINE_WIDTH"]*zoom)+self.squareSide,  (2*BIKE["LINE_WIDTH"]*zoom)+self.squareSide))
        border.fill(BIKE["LINE_COLOUR"])
        bikeSurf = pygame.Surface((self.squareSide, self.squareSide))
        bikeSurf.fill(self.colour)
        border.blit(bikeSurf, (BIKE["LINE_WIDTH"]*zoom, BIKE["LINE_WIDTH"]*zoom))
        border.set_alpha(BIKE["TRANSPARENCY"])
        screen.blit(border, (self.trueX, self.trueY))
        # pygame.draw.rect(screen, self.colour, (self.trueX, self.trueY, self.squareSide, self.squareSide))
        # Draw the agents within the bike
        for index, agent in enumerate(self.agentList.values()):
            # Calculate agent's position within the grid
            row = index // gridSize
            col = index % gridSize
            agentX = self.trueX + agentSize / 2 + (agentSize  * (col)) + (agentPadding * (col + 1))
            agentY = self.trueY + agentSize / 2 + (agentSize  * (row)) + (agentPadding * (row + 1))
            agent.draw(screen, agentX, agentY, zoom)

    def set_agents(self, agentJson:dict, nextRound:dict) -> None:
        """
        Set the agents that are in the bike
        """
        averageEnergy = 0
        averagePedal = 0
        averageSteering = 0
        averageBrake = 0
        averagePoints = 0
        for agentid in self.agentData:
            # Allow for older JSONs that do not have group_id
            if "group_id" not in agentJson[agentid]:
                agentJson[agentid]["group_id"] = 0
            if agentid in nextRound:
                steeringForce = nextRound[agentid]["forces"]["turning"]["steering_force"]
            else:
                steeringForce = 0
            if agentid not in self.agentList:
                self.agentList[agentid] = Agent(self.x, self.y, agentid, agentJson[agentid]["colour"], agentJson[agentid]["group_id"], agentJson[agentid], self.orientation, steeringForce)
            else:
                self.agentList[agentid].update_agent(self.x, self.y, agentid, agentJson[agentid]["colour"], agentJson[agentid]["group_id"], agentJson[agentid], self.orientation, steeringForce)
            # Calculate averages
            averageEnergy += agentJson[agentid]["energy_level"]
            averagePedal += agentJson[agentid]["forces"]["pedal"]
            averageBrake += agentJson[agentid]["forces"]["brake"]
            averageSteering += agentJson[agentid]["forces"]["turning"]["steering_force"]
            averagePoints += agentJson[agentid]["points"]
        if len(self.agentData) == 0:
            averageEnergy = "N/A"
            averagePedal = "N/A"
            averageBrake = "N/A"
            averageSteering = "N/A"
            averagePoints = "N/A"
        else:
            averageEnergy = round(averageEnergy / len(self.agentData), PRECISION)
            averagePedal = round(averagePedal / len(self.agentData), PRECISION)
            averageBrake = round(averageBrake / len(self.agentData), PRECISION)
            averagePoints = round(averagePoints / len(self.agentData), PRECISION)
        avgs = {
            "Average Energy" : averageEnergy,
            "Average Pedal" : averagePedal,
            "Average Steering" : averageSteering,
            "Average Brake" : averageBrake,
            "Average Points" : averagePoints,
        }
        self.properties.update(avgs)

    def get_agents(self) -> dict:
        """
        Return the agents in the bike
        """
        agents = dict()
        for agentid, agent in self.agentList.items():
            agents[agentid] = agent.get_properties()
        return agents

    def propagate_click(self, mouseX:int, mouseY:int, zoom:float) -> None:
        """
        Propagate the click to the agents within the bike
        """
        intersected = False
        for agent in self.agentList.values():
            if agent.click(mouseX, mouseY, zoom):
                intersected = True
        # If an agent was not interacted with, check bike
        if not intersected:
            self.click(mouseX, mouseY, zoom)

    def check_collision(self, mouseX: int, mouseY: int, _:float) -> bool:
        """
        Check if the mouse click intersects with the bike.
        """
        return (self.trueX <= mouseX <= self.trueX + self.squareSide) and \
               (self.trueY <= mouseY <= self.trueY + self.squareSide)

    def draw_overlay(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Overlay agent properties when clicked.
        """
        for _, agent in enumerate(self.agentList.values()):
            agent.set_bike_orientation(self.orientation)
            agent.draw_overlay(screen, offsetX, offsetY, zoom)
        if self.clicked:
            posX = self.nextPos[0] * COORDINATESCALE * zoom + offsetX
            posY = self.nextPos[1] * COORDINATESCALE * zoom + offsetY
            if self.nextPos == (0, 0):
                posX = self.trueX+10 + self.squareSide/2
                posY = self.trueY + self.squareSide/2
            self.draw_arrow(screen, self.colour, (self.trueX + self.squareSide/2, self.trueY + self.squareSide/2), (posX, posY))
            # self.draw_arrow(screen, self.colour, (self.trueX + self.squareSide/2, self.trueY + self.squareSide/2), self.orientation)
        super().draw_overlay(screen, offsetX, offsetY, zoom)
