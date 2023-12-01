"""
Logic for handling bikes in the visualiser
"""
# pylint: disable=import-error, no-name-in-module
import pygame
import pygame_gui
from visualiser.util.Constants import AWDI, OVERLAY
from visualiser.entities.Common import Drawable

class Awdi(Drawable):
    def __init__(self, jsonData:dict) -> None:
        super().__init__("audi", jsonData)
        self.colour = AWDI["COLOUR"]
        properties = {
            "Target" : jsonData["target_bike"],
            "Acceleration" : jsonData["physical_state"]["acceleration"],
            "Velocity" : jsonData["physical_state"]["velocity"],
            "Mass" : jsonData["physical_state"]["mass"],
        }
        self.properties.update(properties)

    def draw(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Draw the lootbox
        """
        # Determine the grid size
        self.trueX = int(self.x * zoom + offsetX)
        self.trueY = int(self.y * zoom + offsetY)
        # Draw the awdi
        border = pygame.Surface(((2*AWDI["LINE_WIDTH"] + AWDI["SIZE"])*zoom, (2*AWDI["LINE_WIDTH"] + AWDI["SIZE"])*zoom))
        border.fill(AWDI["LINE_COLOUR"])
        overlay = pygame.Surface((AWDI["SIZE"]*zoom, AWDI["SIZE"]*zoom))
        overlay.fill(self.colour)
        # Add awdi text
        font = pygame.font.SysFont(OVERLAY["FONT"], int(AWDI["FONT_SIZE"] * zoom))
        text = font.render("owdi", True, "White")
        # center the text
        textX = (AWDI["SIZE"]*zoom - text.get_width()) / 2
        textY = (AWDI["SIZE"]*zoom - text.get_height()) / 2
        overlay.blit(text, (textX, textY))
        # add the overlay to the border
        border.blit(overlay, (AWDI["LINE_WIDTH"]*zoom, AWDI["LINE_WIDTH"]*zoom))
        screen.blit(border, (self.trueX, self.trueY))
        # update the overlay
        self.overlay = self.update_overlay(zoom)
        self.draw_overlay(screen)

    def check_collision(self, mouseX: int, mouseY: int, zoom:float) -> bool:
        """
        Check if the mouse click intersects with the bike.
        """
        return (self.trueX <= mouseX <= self.trueX + AWDI["SIZE"]*zoom) and \
               (self.trueY <= mouseY <= self.trueY + AWDI["SIZE"]*zoom)

    def propagate_click(self, mouseX:int, mouseY:int, zoom:float) -> None:
        """
        Propagate the click
        """
        self.click(mouseX, mouseY, zoom)
