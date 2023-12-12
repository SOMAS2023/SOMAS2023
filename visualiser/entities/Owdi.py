"""
Logic for handling bikes in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
import pygame
import pygame_gui
from visualiser.util.Constants import OWDI, OVERLAY, COORDINATESCALE
from visualiser.entities.Common import Drawable

class Owdi(Drawable):
    def __init__(self, jsonData:dict, targetPosition:dict) -> None:
        self.text = None
        super().__init__("owdi", jsonData)
        self.update_owdi(jsonData, targetPosition)
        self.update_text(1.0)

    def update_owdi(self, jsonData:dict, targetPosition:dict) -> None:
        """
        Update the owdi's properties
        """
        super().update_entity(jsonData)
        if jsonData["target_bike"] in targetPosition:
            location = targetPosition[jsonData["target_bike"]]["physical_state"]["position"]
            self.targetPosition = (location["x"], location["y"])
        else:
            self.targetPosition = (0, 0)
        self.colour = OWDI["COLOUR"]
        properties = {
            "Target" : jsonData["target_bike"],
            "Acceleration" : jsonData["physical_state"]["acceleration"],
            "Velocity" : jsonData["physical_state"]["velocity"],
            "Mass" : jsonData["physical_state"]["mass"],
        }
        self.properties.update(properties)

    def update_text(self, zoom:float) -> None:
        """
        Update the text the agent displays on zoom change
        """
        if zoom == self.zoom and self.text is not None:
            return
        self.zoom = zoom
        font = pygame.font.SysFont(OVERLAY["FONT"], int(OWDI["FONT_SIZE"] * zoom))
        self.text = font.render("owdi", True, "White")

    def draw(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Draw the lootbox
        """
        # Determine the grid size
        self.trueX = int(self.x * COORDINATESCALE * zoom + offsetX - OWDI["SIZE"]*zoom/2)
        self.trueY = int(self.y * COORDINATESCALE * zoom + offsetY - OWDI["SIZE"]*zoom/2)
        # Draw the owdi
        border = pygame.Surface(((2*OWDI["LINE_WIDTH"] + OWDI["SIZE"])*zoom, (2*OWDI["LINE_WIDTH"] + OWDI["SIZE"])*zoom))
        border.fill(OWDI["LINE_COLOUR"])
        overlay = pygame.Surface((OWDI["SIZE"]*zoom, OWDI["SIZE"]*zoom))
        overlay.fill(self.colour)
        # center the text
        self.update_text(zoom)
        textX = (OWDI["SIZE"]*zoom - self.text.get_width()) / 2
        textY = (OWDI["SIZE"]*zoom - self.text.get_height()) / 2
        overlay.blit(self.text, (textX, textY))
        # add the overlay to the border
        border.blit(overlay, (OWDI["LINE_WIDTH"]*zoom, OWDI["LINE_WIDTH"]*zoom))
        # Center the owdi
        screen.blit(border, (self.trueX, self.trueY))

    def check_collision(self, mouseX: int, mouseY: int, zoom:float) -> bool:
        """
        Check if the mouse click intersects with the bike.
        """
        return (self.trueX <= mouseX <= self.trueX + OWDI["SIZE"]*zoom) and \
               (self.trueY <= mouseY <= self.trueY + OWDI["SIZE"]*zoom)

    def propagate_click(self, mouseX:int, mouseY:int, zoom:float) -> None:
        """
        Propagate the click
        """
        self.click(mouseX, mouseY, zoom)

    def draw_overlay(self, screen:pygame_gui.core.UIContainer, offsetX: int, offsetY: int, zoom: float) -> None:
        """
        Draw the overlay for the owdi
        """
        if self.clicked and self.targetPosition != (0, 0):
            self.draw_arrow(screen, self.colour, (self.trueX+OWDI["SIZE"]*zoom/2, self.trueY+OWDI["SIZE"]*zoom/2), (self.targetPosition[0]*COORDINATESCALE*zoom+offsetX, self.targetPosition[1]*COORDINATESCALE*zoom+offsetY))
        super().draw_overlay(screen, offsetX, offsetY, zoom)
