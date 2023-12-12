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
        self.text = None
        super().__init__(lootboxid, jsonData)
        self.colour = COLOURS[jsonData["colour"]]
        self.update_lootbox(jsonData)
        self.update_text(1.0)

    def update_lootbox(self, jsonData:dict) -> None:
        """
        Update the lootbox's properties
        """
        super().update_entity(jsonData)
        properties = {
            "Acceleration" : jsonData["physical_state"]["acceleration"],
            "Velocity" : jsonData["physical_state"]["velocity"],
            "Mass" : jsonData["physical_state"]["mass"],
            "Resources" : jsonData["total_resources"],
            "Colour" : jsonData["colour"].title()
        }
        self.properties.update(properties)
    
    def update_text(self, zoom:float) -> None:
        """
        Update the text the agent displays on zoom change
        """
        if zoom == self.zoom and self.text is not None:
            return
        self.zoom = zoom
        font = pygame.font.SysFont(OVERLAY["FONT"], int(LOOTBOX["FONT_SIZE"] * zoom))
        if self.colour in (COLOURS["white"]):
            self.text = font.render("Lootbox", True, "Black")
        else:
            self.text = font.render("Lootbox", True, "White")

    def draw(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Draw the lootbox
        """
        # Determine the grid size
        self.trueX = int(self.x*COORDINATESCALE * zoom + offsetX - LOOTBOX["WIDTH"]*zoom/2)
        self.trueY = int(self.y*COORDINATESCALE * zoom + offsetY - LOOTBOX["HEIGHT"]*zoom/2)
        # Draw the lootbox
        border = pygame.Surface(((2*LOOTBOX["LINE_WIDTH"] + LOOTBOX["WIDTH"])*zoom, (2*LOOTBOX["LINE_WIDTH"] + LOOTBOX["HEIGHT"])*zoom))
        border.fill(LOOTBOX["LINE_COLOUR"])
        overlay = pygame.Surface((LOOTBOX["WIDTH"]*zoom, LOOTBOX["HEIGHT"]*zoom))
        overlay.fill(self.colour)
        # Add lootbox text
        self.update_text(zoom)
        # center the text
        textX = (LOOTBOX["WIDTH"]*zoom - self.text.get_width()) / 2
        textY = (LOOTBOX["HEIGHT"]*zoom - self.text.get_height()) / 2
        overlay.blit(self.text, (textX, textY))
        # add the overlay to the border
        border.blit(overlay, (LOOTBOX["LINE_WIDTH"]*zoom, LOOTBOX["LINE_WIDTH"]*zoom))
        # Center the lootbox
        screen.blit(border, (self.trueX, self.trueY))

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
