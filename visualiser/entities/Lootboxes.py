"""
Logic for handling bikes in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
import pygame
import pygame_gui
from visualiser.util.Constants import LOOTBOX, OVERLAY, COLOURS, COORDINATESCALE
from visualiser.entities.Common import Drawable

class Lootbox(Drawable):
    def __init__(self, lootboxid:str, jsonData:dict) -> None:
        super().__init__(lootboxid, jsonData)
        self.colour = COLOURS[jsonData["colour"]]
        properties = {
            "Acceleration" : jsonData["physical_state"]["acceleration"],
            "Velocity" : jsonData["physical_state"]["velocity"],
            "Mass" : jsonData["physical_state"]["mass"],
            "Resources" : jsonData["total_resources"],
            "Colour" : jsonData["colour"].title()
        }
        self.properties.update(properties)

    def draw(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Draw the lootbox
        """
        # Determine the grid size
        self.trueX = int(self.x*COORDINATESCALE * zoom + offsetX)
        self.trueY = int(self.y*COORDINATESCALE * zoom + offsetY)
        # Draw the lootbox
        border = pygame.Surface(((2*LOOTBOX["LINE_WIDTH"] + LOOTBOX["WIDTH"])*zoom, (2*LOOTBOX["LINE_WIDTH"] + LOOTBOX["HEIGHT"])*zoom))
        border.fill(LOOTBOX["LINE_COLOUR"])
        overlay = pygame.Surface((LOOTBOX["WIDTH"]*zoom, LOOTBOX["HEIGHT"]*zoom))
        overlay.fill(self.colour)
        # Add lootbox text
        font = pygame.font.SysFont(OVERLAY["FONT"], int(LOOTBOX["FONT_SIZE"] * zoom))
        if self.colour in (COLOURS["white"]):
            text = font.render("Lootbox", True, "Black")
        else:
            text = font.render("Lootbox", True, "White")
        # center the text
        textX = (LOOTBOX["WIDTH"]*zoom - text.get_width()) / 2
        textY = (LOOTBOX["HEIGHT"]*zoom - text.get_height()) / 2
        overlay.blit(text, (textX, textY))
        # add the overlay to the border
        border.blit(overlay, (LOOTBOX["LINE_WIDTH"]*zoom, LOOTBOX["LINE_WIDTH"]*zoom))
        screen.blit(border, (self.trueX, self.trueY))
        # Draw the agents within the bike
        self.overlay = self.update_overlay(zoom)
        self.draw_overlay(screen)

    def check_collision(self, mouseX: int, mouseY: int, zoom:float) -> bool:
        """
        Check if the mouse click intersects with the bike.
        """
        return (self.trueX <= mouseX <= self.trueX + LOOTBOX["WIDTH"]*zoom) and \
               (self.trueY <= mouseY <= self.trueY + LOOTBOX["HEIGHT"]*zoom)

    def propagate_click(self, mouseX:int, mouseY:int, zoom:float) -> None:
        """
        Propagate the click
        """
        self.click(mouseX, mouseY, zoom)
