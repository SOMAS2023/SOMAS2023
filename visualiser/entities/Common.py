"""
Common functions between entities.
"""
import math
import pygame
import pygame_gui
from visualiser.util.Constants import OVERLAY, COORDINATESCALE, PRECISION, DIM, ARROWS, COLOURS
class Drawable:
    def __init__(self, entityid:str, jsonData:dict, x=None, y=None) -> None:
        if x is None or y is None:
            self.x = jsonData["physical_state"]["position"]["x"]
            self.y = jsonData["physical_state"]["position"]["y"]
        else:
            self.x = x
            self.y = y
        self.x = round(self.x, PRECISION)
        self.y = round(self.y, PRECISION)
        self.trueX = round(self.x*COORDINATESCALE, PRECISION)
        self.trueY = round(self.y*COORDINATESCALE, PRECISION)
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
        return border

    def check_collision(self, mouseX:int, mouseY:int, zoom:float) -> bool:
        """
        Check if the mouse click intersects with the agent.
        """
        raise NotImplementedError

    def draw_overlay(self, screen:pygame_gui.core.UIContainer, offsetX:int, offsetY:int, zoom:float) -> None:
        """
        Overlay agent properties when clicked.
        """
        if self.clicked:
            screen.blit(self.overlay, (self.trueX, self.trueY))

    def get_properties(self) -> dict:
        """
        Return the properties of the agent.
        """
        properties = {
            "X" : self.x,
            "Y" : self.y,
        }
        properties.update(self.properties)
        return properties

    def draw_arrow(self, screen, colour, startPoint, secondArg):
        """
        Draw a line with multiple arrowheads to indicate direction.
        """
        arrowLength = ARROWS["ARROW_LENGTH"]
        arrowAngle = ARROWS["ARROW_ANGLE"] * math.pi / 180
        numArrows = ARROWS["NUM_ARROWS"]
        length = 10
        if colour == COLOURS["white"]:
            colour = "black"
        # Determine if secondArg is an endPoint or an orientation
        if isinstance(secondArg, tuple):  # If secondArg is a tuple, assume it's an endPoint
            endPoint = secondArg
        else:  # If secondArg is not a tuple, calculate endPoint based on orientation
            orientation = secondArg
            angle = orientation * math.pi
            endPoint = (startPoint[0] + length * math.cos(angle), startPoint[1] + length * math.sin(angle))
        # Calculate the direction of the line
        dx, dy = endPoint[0] - startPoint[0], endPoint[1] - startPoint[1]
        gradient = dy / dx if dx != 0 else float('inf')

        # Place arrowheads along the line
        def get_intersections(x, y, grad):
            points = []
            if grad != 0:
                yLeft = y - grad * x
                if 0 <= yLeft <= DIM["GAME_SCREEN_HEIGHT"]:
                    points.append((0, yLeft))
                yRight = y + grad * (DIM["GAME_SCREEN_WIDTH"] - x)
                if 0 <= yRight <= DIM["GAME_SCREEN_HEIGHT"]:
                    points.append((DIM["GAME_SCREEN_WIDTH"], yRight))
                xTop = x - y / grad
                if 0 <= xTop <= DIM["GAME_SCREEN_WIDTH"]:
                    points.append((xTop, 0))
                xBottom = x + (DIM["GAME_SCREEN_HEIGHT"] - y) / grad
                if 0 <= xBottom <= DIM["GAME_SCREEN_WIDTH"]:
                    points.append((xBottom, DIM["GAME_SCREEN_HEIGHT"]))
            else:
                points.extend([(0, y), (DIM["GAME_SCREEN_WIDTH"], y)])
            return points

        def draw_arrowhead(screen, color, tip, direction, length, angle):
            dx1, dy1 = length * math.cos(direction + angle), length * math.sin(direction + angle)
            dx2, dy2 = length * math.cos(direction - angle), length * math.sin(direction - angle)
            pygame.draw.polygon(screen, color, [tip, (tip[0] - dx1, tip[1] - dy1), (tip[0] - dx2, tip[1] - dy2)])

        # Get intersections with screen boundaries
        intersections = get_intersections(startPoint[0], startPoint[1], gradient)
        if len(intersections) == 2:
            pygame.draw.line(screen, colour, intersections[0], intersections[1], 3)
            lineDir = math.atan2(dy, dx)
            for i in range(1, numArrows + 1):
                fraction = i / (numArrows + 1)
                x, y = [p0 + fraction * (p1 - p0) for p0, p1 in zip(intersections[0], intersections[1])]
                draw_arrowhead(screen, colour, (x, y), lineDir, arrowLength, arrowAngle)
