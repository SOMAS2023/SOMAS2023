"""
Logic for the game screen visualiser
"""
# pylint: disable=no-member, import-error, no-name-in-module
import pygame
import pygame_gui
from pygame_gui.elements import UIButton, UILabel
from visualiser.util.Constants import DIM, BGCOLOURS, MAXZOOM, MINZOOM
from visualiser.util.HelperFunc import make_center
from visualiser.entities.Bikes import Bike

class GameScreen:
    def __init__(self) -> None:
        self.UIElements = {}
        self.round = 0
        self.offsetX = 0
        self.offsetY = 0
        self.dragging = False
        self.mouseX, self.mouseY = 0, 0
        self.oldZoom = 1.0
        self.zoom = 1.0
        self.bikes = {'0' : Bike(100, 100, 0)}
        self.elements = {}

    def init_ui(self, manager:pygame_gui.UIManager, screen:pygame_gui.core.UIContainer) -> dict:
        """
        Initialise the UI for the main menu screen
        """
        x, _ = make_center((DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"]), (DIM["UI_WIDTH"], DIM["SCREEN_HEIGHT"]))
        self.elements["reset"] = UIButton(
            relative_rect=pygame.Rect((x, 10), (DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"])),
            text="Reset",
            manager=manager,
            container=screen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        topmargin = 250
        # Round count
        self.elements["round_count"] = UILabel(
            relative_rect=pygame.Rect((x, topmargin+DIM["BUTTON_HEIGHT"]), (DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"])),
            text="Round: 0",
            manager=manager,
            container=screen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        # Round controls
        factor = 0.85
        x, _ = make_center((DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"]), (DIM["UI_WIDTH"], DIM["SCREEN_HEIGHT"]))
        self.elements["increase_round"] = UIButton(
            relative_rect=pygame.Rect((x, topmargin), (DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"])),
            text="Increase Round",
            manager=manager,
            container=screen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        self.elements["decrease_round"] = UIButton(
            relative_rect=pygame.Rect((x, topmargin+2*DIM["BUTTON_HEIGHT"]), (DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"])),
            text="Decrease Round",
            manager=manager,
            container=screen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        return self.elements

    def render_game_screen(self, screen:pygame_gui.core.UIContainer) -> None:
        """
        Render the game screen
        """
        screen.fill((255, 255, 255))
        self.draw_grid(screen)
        # Draw agents
        self.bikes["0"].draw(screen, self.offsetX, self.offsetY, self.zoom)
        # Divider line
        lineWidth = 1
        pygame.draw.line(screen, "#555555", (DIM["GAME_SCREEN_WIDTH"]-lineWidth, 0), (DIM["GAME_SCREEN_WIDTH"]-lineWidth, DIM["SCREEN_HEIGHT"]), lineWidth)

    def process_events(self, event:pygame.event.Event) -> None:
        """
        Process events in the game screen
        """
        match event.type:
            case pygame.MOUSEBUTTONDOWN:
                match event.button:
                    case 1:  # Left click
                        self.dragging = True
                        self.mouseX, self.mouseY = event.pos
                        self.propagate_click(event.pos)
                    case 4:  # Scroll up
                        self.adjust_zoom(1.1, event.pos)
                    case 5:  # Scroll down
                        self.adjust_zoom(0.9, event.pos)
            case pygame.MOUSEBUTTONUP:
                if event.button == 1:  # Left click
                    self.propagate_click((-1, -1))
                    self.dragging = False
            case pygame.MOUSEMOTION:
                if self.dragging:
                    mouseX, mouseY = event.pos
                    # Reverse the direction of the offset updates
                    self.offsetX += mouseX - self.mouseX
                    self.offsetY += mouseY - self.mouseY
                    self.mouseX, self.mouseY = mouseX, mouseY
            case pygame_gui.UI_BUTTON_PRESSED:
                if event.ui_element == self.elements["increase_round"]:
                    self.change_round(self.round + 1)
                elif event.ui_element == self.elements["decrease_round"]:
                    self.change_round(self.round - 1)

    def change_round(self, newRound:int) -> None:
        """
        Change the current round
        """
        self.round = max(0, newRound)
        self.elements["round_count"].set_text(f"Round: {self.round}")

    def propagate_click(self, mousePos:tuple) -> None:
        """
        Propagate the click to all entities
        """
        mouseX, mouseY = mousePos
        self.bikes["0"].propagate_click(mouseX, mouseY, self.offsetX, self.offsetY, self.zoom)

    def adjust_zoom(self, zoomFactor:float, mousePos:tuple) -> None:
        """
        Adjust the zoom level of the game screen
        """
        mouseX, mouseY = mousePos
        self.oldZoom = self.zoom
        self.zoom = max(MINZOOM, min(self.zoom * zoomFactor, MAXZOOM))
        # Calculate the new offsets
        self.offsetX = mouseX - (mouseX - self.offsetX) * (self.zoom / self.oldZoom)
        self.offsetY = mouseY - (mouseY - self.offsetY) * (self.zoom / self.oldZoom)

    def draw_grid(self, surface:pygame.Surface) -> None:
        """
        Draw the grid on the game screen
        """
        zoomedSpacing = DIM["GRIDSPACING"] * self.zoom

        # Correctly applying the offsets
        startX = self.offsetX % zoomedSpacing
        startY = self.offsetY % zoomedSpacing
        width = DIM["GAME_SCREEN_WIDTH"]
        height = DIM["SCREEN_HEIGHT"]
        # Draw vertical lines
        for x in range(-int(zoomedSpacing) + int(startX), width, int(zoomedSpacing)):
            pygame.draw.line(surface, BGCOLOURS["GRID"], (x, 0), (x, height))

        # Draw horizontal lines
        for y in range(-int(zoomedSpacing) + int(startY), height, int(zoomedSpacing)):
            pygame.draw.line(surface, BGCOLOURS["GRID"], (0, y), (width, y))
