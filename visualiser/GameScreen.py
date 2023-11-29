"""
Logic for the game screen visualiser
"""
# pylint: disable=no-member, import-error, no-name-in-module
import random
import colorsys
import pygame
import pygame_gui
from pygame_gui.elements import UIButton, UILabel, ui_text_box
from visualiser.util.Constants import DIM, BGCOLOURS, MAXZOOM, MINZOOM, ZOOM, BIKE, OVERLAY, PRECISION
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
        self.manager = None
        self.playSpeed = 1
        self.isPlaying = False
        self.playEvent = pygame.USEREVENT + 100
        self.mouseXCur = 0
        self.mouseYCur = 0

    def init_ui(self, manager:pygame_gui.UIManager, screen:pygame_gui.core.UIContainer) -> dict:
        """
        Initialise the UI for the main menu screen
        """
        self.manager = manager
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
        #control information
        self.elements["controls"] = ui_text_box.UITextBox(
            relative_rect=pygame.Rect((x, 10+DIM["BUTTON_HEIGHT"]), (DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"]*3.5)),
            html_text="<font face=verdana size=3 color=#FFFFFF><b>Controls</b></font><br><br><font face=verdana size=3 color=#FFFFFF><b>Space</b> - Play/Pause<br><b>Right</b> - Next Round<br><b>Left</b> - Previous Round<br><b>Up</b> - Increase Speed<br><b>Down</b> - Decrease Speed<br><b>Scroll</b> - Zoom<br><b>Click</b> - Select Entity</font>",
            manager=manager,
            container=screen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            },
            object_id="#controls"
        )
        topmargin = 325
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
        self.elements["round_slider"] = pygame_gui.elements.UIHorizontalSlider(
            relative_rect=pygame.Rect((x, topmargin+DIM["BUTTON_HEIGHT"]*2), (DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"]//2)),
            start_value=0,
            value_range=(0, self.maxRound),
            click_increment=self.maxRound//10,
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
        factor = 1
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
            relative_rect=pygame.Rect((x, topmargin+2.5*DIM["BUTTON_HEIGHT"]), (DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"])),
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
        # play pause button
        self.elements["play_pause"] = UIButton(
            relative_rect=pygame.Rect((x, topmargin+DIM["BUTTON_HEIGHT"]*4), (DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"]//2)),
            text="Play",
            manager=manager,
            container=screen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        # play pause speed
        self.elements["play_pause_speed"] = UILabel(
            relative_rect=pygame.Rect((x, topmargin+DIM["BUTTON_HEIGHT"]*4.5), (DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"]//2)),
            text="1 Round/Sec",
            manager=manager,
            container=screen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        # play pause speed
        self.elements["play_pause_speed_slider"] = pygame_gui.elements.UIHorizontalSlider(
            relative_rect=pygame.Rect((x, topmargin+DIM["BUTTON_HEIGHT"]*5), (DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"]//2)),
            start_value=1,
            value_range=(1, 10),
            click_increment=1,
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
        # Draw awdi
        self.awdi.draw(screen, self.offsetX, self.offsetY, self.zoom)
        # # Draw lootboxes
        for lootbox in self.lootboxes:
            lootbox.draw(screen, self.offsetX, self.offsetY, self.zoom)
        # Draw agents
        for bike in self.bikes:
            bike.draw(screen, self.offsetX, self.offsetY, self.zoom)
        self.draw_mouse_coords(screen)
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
                    # Mouse inputs
                    case 1:  # Left click
                        if event.pos < (DIM["GAME_SCREEN_WIDTH"], DIM["SCREEN_HEIGHT"]):
                            self.dragging = True
                        self.mouseX, self.mouseY = event.pos
                        self.propagate_click(event.pos)
                    case 4:  # Scroll up
                        self.adjust_zoom(1.1, event.pos)
                    case 5:  # Scroll down
                        self.adjust_zoom(0.9, event.pos)
            #Interact with UI
            case pygame.MOUSEBUTTONUP:
                if event.button == 1:  # Left click
                    self.propagate_click((-1, -1))
                    self.dragging = False
            #Pan screen
            case pygame.MOUSEMOTION:
                self.mouseXCur, self.mouseYCur = event.pos
                if self.dragging:
                    # Reverse the direction of the offset updates
                    self.offsetX += self.mouseXCur - self.mouseX
                    self.offsetY += self.mouseYCur - self.mouseY
                    self.mouseX, self.mouseY = self.mouseXCur, self.mouseYCur
            # Handle key presses
            case pygame.KEYDOWN:
                match event.key:
                    # Space bar pause/plays
                    case pygame.K_SPACE:
                        self.toggle_play()
                    # Arrow keys advance rounds
                    case pygame.K_RIGHT:
                        self.change_round(self.round + 1)
                    case pygame.K_LEFT:
                        self.change_round(self.round - 1)
                    # Up down increases speed
                    case pygame.K_UP:
                        self.change_speed(self.playSpeed + 1)
                    case pygame.K_DOWN:
                        self.change_speed(self.playSpeed - 1)
            # Handle UI buttons
            case pygame_gui.UI_BUTTON_PRESSED:
                if event.ui_element == self.elements["increase_round"]:
                    self.change_round(self.round + 1)
                elif event.ui_element == self.elements["decrease_round"]:
                    self.change_round(self.round - 1)
                # Play pause
                elif event.ui_element == self.elements["play_pause"]:
                    self.toggle_play()
            #Sliders
            case pygame_gui.UI_HORIZONTAL_SLIDER_MOVED:
                if event.ui_element == self.elements["round_slider"]:
                    self.change_round(event.value)
                elif event.ui_element == self.elements["play_pause_speed_slider"]:
                    self.change_speed(event.value)
            case self.playEvent:
                self.elements["round_slider"].set_current_value(self.round + 1)
                self.change_round(self.round + 1)

    def toggle_play(self) -> None:
        """
        Toggle the play/pause button
        """
        if self.isPlaying:
            self.elements["play_pause"].set_text("Play")
            self.isPlaying = False
            pygame.time.set_timer(self.playEvent, 0)
        else:
            self.elements["play_pause"].set_text("Pause")
            self.isPlaying = True
            pygame.time.set_timer(self.playEvent, int(1000//self.playSpeed))

    def change_speed(self, newSpeed:int) -> None:
        """
        Change the speed of the game
        """
        self.playSpeed = min(10, max(1, newSpeed))
        self.elements["play_pause_speed_slider"].set_current_value(self.playSpeed)
        self.elements["play_pause_speed"].set_text(f"{self.playSpeed} Round/Sec")
        if self.isPlaying:
            pygame.time.set_timer(self.playEvent, int(1000//self.playSpeed))

    def change_round(self, newRound:int) -> None:
        """
        Change the current round using the round controls and json data
        """
        self.round = max(0, min(self.maxRound, newRound))
        self.elements["round_count"].set_text(f"Round: {self.round}")
        #Reload bikes
        self.bikes = []
        for bikeid, bike in self.jsonData[self.round]["bikes"].items():
            if bikeid not in self.bikeColourMap:
                self.bikeColourMap[bikeid] = self.allocate_colour()
            self.bikes.append(Bike(bikeid, bike, self.bikeColourMap[bikeid], self.jsonData[self.round]["agents"]))
        # Reload lootboxes
        self.lootboxes = []
        for lootboxid, lootbox in self.jsonData[self.round]["loot_boxes"].items():
            self.lootboxes.append(Lootbox(lootboxid, lootbox))
        self.awdi = Awdi(self.jsonData[self.round]["audi"])

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
            bike.propagate_click(mouseX, mouseY, self.zoom)
        for lootbox in self.lootboxes:
            lootbox.propagate_click(mouseX, mouseY, self.zoom)
        self.awdi.propagate_click(mouseX, mouseY, self.zoom)

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

    def draw_mouse_coords(self, surface:pygame.Surface) -> None:
        """
        Draw the mouse coordinates on the game screen
        """
        font = pygame.font.SysFont("Arial", 15)
        x = round(self.mouseXCur / self.zoom - self.offsetX / self.zoom, PRECISION)
        y = round(self.mouseYCur / self.zoom - self.offsetY / self.zoom, PRECISION)
        text = font.render(f"({x}, {y})", True, (0, 0, 0))
        surface.blit(text, (OVERLAY["FPS_PAD"], OVERLAY["FPS_PAD"]))

    def set_json(self, data:dict) -> None:
        """
        Set the data for the game screen
        """
        if data is None:
            return
        self.jsonData = data
        self.maxRound = len(data) - 1
