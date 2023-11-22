"""
Logic for the game screen visualiser
"""
# pylint: disable=no-member, import-error, no-name-in-module
import random
import colorsys
import pygame
import pygame_gui
from pygame_gui.elements import UIButton, UILabel
from visualiser.util.Constants import DIM, BGCOLOURS, MAXZOOM, MINZOOM, ZOOM, COORDINATESCALE, BIKE
from visualiser.util.HelperFunc import make_center
from visualiser.entities.Bikes import Bike
from visualiser.entities.Lootboxes import Lootbox
from visualiser.entities.Awdi import Awdi

class GameScreen:
    def __init__(self) -> None:
        self.UIElements = {}
        self.round = 0
        self.offsetX = 0
        self.offsetY = 0
        self.dragging = False
        self.mouseX, self.mouseY = 0, 0
        self.oldZoom = 1.0
        self.zoom = ZOOM
        self.bikes = []
        self.lootboxes = []
        self.awdi = None
        self.elements = {}
        self.jsonData = None
        self.maxRound = 0
        self.bikeColourMap = {}

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
        # Draw lootboxes
        for lootbox in self.lootboxes:
            lootbox.draw(screen, self.offsetX, self.offsetY, self.zoom)
        # # Draw agents
        for bike in self.bikes:
            bike.draw(screen, self.offsetX, self.offsetY, self.zoom)
        # Draw awdi
        self.awdi.draw(screen, self.offsetX, self.offsetY, self.zoom)
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
        self.round = max(0, min(self.maxRound, newRound))
        self.elements["round_count"].set_text(f"Round: {self.round}")
        #Reload bikes
        self.bikes = []
        data = self.jsonData[f"loop_{self.round}"]["bikes"]
        for b in data:
            if data[b]["id"] not in self.bikeColourMap:
                self.bikeColourMap[data[b]["id"]] = self.allocate_colour()
            self.bikes.append(Bike(data[b]["position"]["x"]*COORDINATESCALE, data[b]["position"]["y"]*COORDINATESCALE, data[b]["id"], self.bikeColourMap[data[b]["id"]]))
        # Update the agents and bikes
        for b in self.bikes:
            b.change_round(self.jsonData[f"loop_{self.round}"]["bikes"])
        # Reload lootboxes
        self.lootboxes = []
        lootboxes = self.jsonData[f"loop_{self.round}"]["lootboxes"]
        for l in lootboxes:
            self.lootboxes.append(Lootbox(lootboxes[l]["position"]["x"]*COORDINATESCALE, lootboxes[l]["position"]["y"]*COORDINATESCALE, l))
        # Update the lootboxes
        for l in self.lootboxes:
            l.change_round(self.jsonData[f"loop_{self.round}"]["lootboxes"])
        self.awdi = Awdi(self.jsonData[f"loop_{self.round}"]["awdi"]["position"]["x"]*COORDINATESCALE, \
                         self.jsonData[f"loop_{self.round}"]["awdi"]["position"]["y"]*COORDINATESCALE, self.jsonData[f"loop_{self.round}"]["awdi"]["id"])
        self.awdi.change_round(self.jsonData[f"loop_{self.round}"]["awdi"])

    def allocate_colour(self) -> str:
        """
        Allocate a colour to a bike
        """
        hue = random.randint(BIKE["COLOURS"]["MINHUE"], BIKE["COLOURS"]["MAXHUE"]) / 360
        saturation = random.randint(BIKE["COLOURS"]["MINSAT"], BIKE["COLOURS"]["MAXSAT"]) / 100
        value = random.randint(BIKE["COLOURS"]["MINVAL"], BIKE["COLOURS"]["MAXVAL"]) / 100
        colour = colorsys.hsv_to_rgb(hue, saturation, value)
        colour = (colour[0] * 255, colour[1] * 255, colour[2] * 255)
        return colour

    def propagate_click(self, mousePos:tuple) -> None:
        """
        Propagate the click to all entities
        """
        mouseX, mouseY = mousePos
        for bike in self.bikes:
            bike.propagate_click(mouseX, mouseY, self.offsetX, self.offsetY, self.zoom)
        for lootbox in self.lootboxes:
            lootbox.propagate_click(mouseX, mouseY, self.offsetX, self.offsetY, self.zoom)
        self.awdi.propagate_click(mouseX, mouseY, self.offsetX, self.offsetY, self.zoom)

    def adjust_zoom(self, zoomFactor:float, mousePos:tuple) -> None:
        """
        Adjust the zoom level of the game screen
        """
        mouseX, mouseY = mousePos
        self.oldZoom = self.zoom
        self.zoom = max(MINZOOM, min(self.zoom * zoomFactor, MAXZOOM))
        # Calculate the new offsets
        self.offsetX = mouseX - (mouseX - self.offsetX) * (self.zoom / self.oldZoom)
        self.offsetY = mouseY  - (mouseY  - self.offsetY) * (self.zoom / self.oldZoom)

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

    def set_data(self, data:dict) -> None:
        """
        Set the data for the game screen
        """
        self.jsonData = data
        self.maxRound = len(data) - 1
