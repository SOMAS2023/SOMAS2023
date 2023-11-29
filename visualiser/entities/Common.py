"""
Common functions between entities.
"""
import pygame
import pygame_gui
from visualiser.util.Constants import OVERLAY, COORDINATESCALE, PRECISION
class Drawable:
    def __init__(self, entityid:str, jsonData:dict, x=None, y=None) -> None:
        if x is None or y is None:
            self.x = round(jsonData["physical_state"]["position"]["x"]*COORDINATESCALE, PRECISION)
            self.y = round(jsonData["physical_state"]["position"]["y"]*COORDINATESCALE, PRECISION)
        else:
            self.x = round(x*COORDINATESCALE, PRECISION)
            self.y = round(y*COORDINATESCALE, PRECISION)
        self.trueX = self.x
        self.trueY = self.y
        self.clicked = False
        self.id = entityid
        self.properties = {
            "ID" : self.id,
            "Position" : f"{self.x}, {self.y}",
        }
        self.overlay = pygame.Surface((0, 0))

    def click(self, mouseX:int, mouseY:int, zoom:float) -> bool:
        """
        Check if the mouse click intersects with the agent.
        """
        if (mouseX, mouseY) == (-1, -1):
            self.clicked = False
            return False
        self.clicked = self.check_collision(mouseX, mouseY, zoom)
        if self.clicked:
            print(f"Clicked on {self.id} at ({self.x}, {self.y})")
        return self.clicked

    def update_overlay(self, zoom:float) -> pygame_gui.core.UIContainer:
        """
        Overlay entity properties when clicked.
        Returns True if the entity was clicked on.
        """
        zoom = 1.0 # disable zoom
        # Basic setup
        font = pygame.font.SysFont(OVERLAY["FONT"], int(OVERLAY["FONT_SIZE"] * zoom))
        lineSpacing = OVERLAY["PADDING"] * zoom  # Apply zoom to line spacing

        # Calculate the total height required for the overlay
        textHeight = font.size("Test")[1]  # Height of one line of text
        totalHeight = len(self.properties) * (textHeight + lineSpacing)

        # Create the overlay surface with the calculated height
        overlay = pygame.Surface((OVERLAY["WIDTH"] * zoom, totalHeight))
        overlay.fill(OVERLAY["BACKGROUND_COLOUR"])
        length, height = overlay.get_size()

         # Draw each property on left and right edge
        for i, (attr, value) in enumerate(self.properties.items()):
            # Draw attribute name (left-aligned)
            text = font.render(attr, True, OVERLAY["TEXT_COLOUR"])
            overlay.blit(text, (OVERLAY["PADDING"] * zoom, i * (lineSpacing + textHeight)))

            # Draw attribute value (right-aligned)
            valueText = font.render(str(value), True, OVERLAY["TEXT_COLOUR"])
            valueTextX = length - valueText.get_width() - OVERLAY["PADDING"] * zoom
            overlay.blit(valueText, (valueTextX, i * (lineSpacing + textHeight)))

            # Draw dividing line only if it's not the last item
            if i < len(self.properties) - 1:
                lineY = (i + 1) * (lineSpacing + textHeight)
                pygame.draw.line(overlay, OVERLAY["LINE_COLOUR"], (0, lineY), (length, lineY), OVERLAY["LINE_WIDTH"])
        border = pygame.Surface((length + 2 * OVERLAY["BORDER_WIDTH"], height + 2 * OVERLAY["BORDER_WIDTH"]))
        border.fill(OVERLAY["BORDER_COLOUR"])
        border.blit(overlay, (OVERLAY["BORDER_WIDTH"], OVERLAY["BORDER_WIDTH"]))

        # Make transparent
        border.set_alpha(OVERLAY["TRANSPARENCY"])
        return border

    def check_collision(self, mouseX:int, mouseY:int, zoom:float) -> bool:
        """
        Check if the mouse click intersects with the agent.
        """
        raise NotImplementedError

    def draw_overlay(self, screen:pygame_gui.core.UIContainer) -> None:
        """
        Overlay agent properties when clicked.
        """
        if self.clicked:
            screen.blit(self.overlay, (self.trueX, self.trueY))
