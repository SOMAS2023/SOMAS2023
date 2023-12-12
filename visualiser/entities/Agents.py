"""
Logic for handling the agents in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
import math
import pygame
import pygame_gui
from visualiser.util.Constants import AGENT, COLOURS, PRECISION
from visualiser.entities.Common import Drawable

class Agent(Drawable):
    def __init__(self, x:int, y:int, agentid:str, colour:pygame.color, groupID, jsonData:dict, bikeOrientation, nextOrient) -> None:
        self.bikeOrientation = bikeOrientation
        self.text = None
        super().__init__(agentid, jsonData, x, y)
        self.update_agent(x, y, colour, groupID, jsonData, bikeOrientation, nextOrient)
        self.update_text(1.0)

    def update_agent(self, x:int, y:int, colour:pygame.color, groupID, jsonData:dict, bikeOrientation, nextOrient) -> None:
        """
        Update the agent's properties
        """
        super().update_entity(jsonData, x, y)
        self.colour = COLOURS[colour]
        self.radius = AGENT["SIZE"]
        self.onBike = jsonData["on_bike"]
        self.bikeOrientation = bikeOrientation
        self.nextOrient = nextOrient
        self.steeringForce = round(jsonData["forces"]["turning"]["steering_force"], PRECISION)
        if groupID == 0:
            self.groupID = "-"
        else:
            self.groupID = str(groupID)
        properties = {
            "Pedal" : jsonData["forces"]["pedal"],
            "Brake" : jsonData["forces"]["brake"],
            "Colour" : colour.title(),
            "Steering?" : f'{jsonData["forces"]["turning"]["steer_bike"] != 0}, {self.steeringForce}',
            "Energy" : round(jsonData["energy_level"], PRECISION),
            "Points" : jsonData["points"],
            "GroupID" : self.groupID,
        }
        self.properties.update(properties)

    def update_text(self, zoom:float) -> None:
        """
        Update the text the agent displays on zoom change
        """
        if zoom == self.zoom and self.text is not None:
            return
        self.zoom = zoom
        font = pygame.font.SysFont("Arial", int(AGENT["FONT_SIZE"]*zoom))
        if self.colour in (COLOURS["white"]):
            self.text = font.render(self.groupID, True, (0, 0, 0))
        else:
            self.text = font.render(self.groupID, True, (255, 255, 255))


    def check_collision(self, mouseX:int, mouseY:int, zoom:float) -> bool:
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
        self.trueX = offsetX
        self.trueY = offsetY
        pygame.draw.circle(screen, AGENT["LINE_COLOUR"], (self.trueX, self.trueY), zoomedRadius+max(AGENT["LINE_WIDTH"]*zoom, 1))
        pygame.draw.circle(screen, self.colour, (self.trueX, self.trueY), zoomedRadius)
        # Draw group ID
        self.update_text(zoom)
        textRect = self.text.get_rect()
        textRect.center = (self.trueX, self.trueY)
        screen.blit(self.text, textRect)

    def get_properties(self) -> dict:
        """
        Return the properties of the agent.
        """
        properties =  super().get_properties()
        properties["onBike"] = self.onBike
        return properties

    def set_bike_orientation(self, orientation:float) -> None:
        """
        Set the orientation of the bike.
        """
        self.bikeOrientation = orientation

    def draw_overlay(self, screen: pygame_gui.core.UIContainer, offsetX: int, offsetY: int, zoom: float) -> None:
        """
        Overlay agent properties when clicked.
        """
        if self.clicked:
            angle = self.bikeOrientation+self.steeringForce
            if angle > 1:
                angle -= 2
            elif angle < -1:
                angle += 2
            # self.draw_arrow(screen, self.colour, (self.trueX, self.trueY), angle)
        super().draw_overlay(screen, offsetX, offsetY, zoom)
