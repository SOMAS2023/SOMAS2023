"""
Common functions between entities.
"""
import pygame
import pygame_gui
class Drawable:
    def __init__(self, x, y) -> None:
        self.x = x
        self.y = y
        self.trueX = x
        self.trueY = y
        self.properties = {"test" : 41}
        self.clicked = False
        self.id = None
        self.overlay = pygame.Surface((0, 0))

    def update_position(self, x:int, y:int) -> None:
        """
        Update the position of the agent.
        """
        self.x = x
        self.y = y

    def click(self, mouseX:int, mouseY:int, offsetX:int, offsetY:int, zoom:float) -> bool:
        """
        Check if the mouse click intersects with the agent.
        """
        if (mouseX, mouseY) == (-1, -1):
            self.clicked = False
            return False
        self.clicked = self.check_collision(mouseX, mouseY, offsetX, offsetY, zoom)
        if self.clicked:
            print(f"Clicked on {self.id} at ({self.x}, {self.y})")
        return self.clicked

    def update_overlay(self, zoom:float) -> pygame_gui.core.UIContainer:
        """
        Overlay entity properties when clicked.
        Returns True if the entity was clicked on.
        """
        # Basic setup
        font = pygame.font.SysFont("Arial", int(20 * zoom))  # Adjust the font size based on zoom
        padding = 10  # Space between rows
        overlay = pygame.Surface((200, len(self.properties) * (20 * zoom + padding)))  # Adjust size as needed
        overlay.fill((255, 255, 255))  # Background color of the overlay, change as needed

        # Draw each property
        for i, (attr, value) in enumerate(self.properties.items()):
            text = font.render(f"{attr}: {value}", True, (0, 0, 0))  # Text color
            overlay.blit(text, (10, i * (20 * zoom + padding)))
        return overlay

    def check_collision(self, mouseX:int, mouseY:int, offsetX:int, offsetY:int, zoom:float) -> bool:
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
