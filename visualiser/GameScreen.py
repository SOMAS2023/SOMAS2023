"""
Logic for the game screen visualiser
"""
# pylint: disable=no-member, import-error, no-name-in-module
import pygame
import pygame_gui
from pygame_gui.elements import UIButton, UILabel, ui_text_box
from pygame_gui.core import UIContainer
from visualiser.util.Constants import DIM, BGCOLOURS, MAXZOOM, MINZOOM, ZOOM, FPSDISPLAYRATE, OVERLAY, PRECISION, CONSOLE, EPSILON, COORDINATESCALE, ENERGYTHRESHOLD, MAXSPEED, ITERATIONLENGTH
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
        self.bikes = {}
        self.lootboxes = {}
        self.awdi = None
        self.elements = {}
        self.jsonData = None
        self.maxRound = 0
        self.bikeColourMap = {}
        self.manager = None
        self.playSpeed = 1
        self.isPlaying = False
        self.playEvent = pygame.USEREVENT + 100
        self.mouseDraw = pygame.USEREVENT + 102
        self.mouseSurface = None
        pygame.time.set_timer(self.mouseDraw, int(1000/FPSDISPLAYRATE))
        self.mouseXCur = 0
        self.mouseYCur = 0
        self.consoleScreen = None
        self.aliveAgents = 0
        self.stats = {
            "Active Bikes" : 0,
            "Active Lootboxes" : 0,
            "Alive Agents" : 0,
            "Dead Agents" : 0,
        }
        self.agents = {}
        self.maxDead = 0
        self.draw_mouse_coords()

    def init_ui(self, manager:pygame_gui.UIManager, uiscreen:UIContainer, consoleContainer:UIContainer) -> dict:
        """
        Initialise the UI for the main menu screen
        """
        self.manager = manager
        x, _ = make_center((DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"]), (DIM["UI_WIDTH"], DIM["SCREEN_HEIGHT"]))
        self.elements["reset"] = UIButton(
            relative_rect=pygame.Rect((x, 10), (DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"])),
            text="Reset",
            manager=manager,
            container=uiscreen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        #control information
        self.elements["controls"] = ui_text_box.UITextBox(
            relative_rect=pygame.Rect((x, 10+DIM["BUTTON_HEIGHT"]), (DIM["BUTTON_WIDTH"], 220)),
            html_text="<font face=verdana size=3 color=#FFFFFF><b>Controls</b></font><br><font face=verdana size=3 color=#FFFFFF><b>Space</b> - Play/Pause<br><b>Right</b> - Next Iteration<br><b>Left</b> - Previous Iteration<br><b>Up</b> - Increase Speed<br><b>Down</b> - Decrease Speed<br><b>Scroll</b> - Zoom<br><b>Click</b> - Select Entity</font>", # pylint: disable=line-too-long
            manager=manager,
            container=uiscreen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            },
            object_id="#controls"
        )
        topmargin = 300
        # Round count
        self.elements["round_count"] = UILabel(
            relative_rect=pygame.Rect((x, topmargin+DIM["BUTTON_HEIGHT"]), (DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"])),
            text="Iteration: 0",
            manager=manager,
            container=uiscreen,
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
            click_increment=ITERATIONLENGTH,
            manager=manager,
            container=uiscreen,
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
            text="Increase Iteration",
            manager=manager,
            container=uiscreen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        # Iteration count
        self.elements["iteration_count"] = UILabel(
            relative_rect=pygame.Rect((x, topmargin+DIM["BUTTON_HEIGHT"]*2.5), (DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"]/2)),
            text="Round: 0",
            manager=manager,
            container=uiscreen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        self.elements["decrease_round"] = UIButton(
            relative_rect=pygame.Rect((x, topmargin+3*DIM["BUTTON_HEIGHT"]), (DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"])),
            text="Decrease Iteration",
            manager=manager,
            container=uiscreen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        topmargin += 15
        # play pause button
        self.elements["play_pause"] = UIButton(
            relative_rect=pygame.Rect((x, topmargin+DIM["BUTTON_HEIGHT"]*4), (DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"]//2)),
            text="Play",
            manager=manager,
            container=uiscreen,
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
            text="1 Iteration/Sec",
            manager=manager,
            container=uiscreen,
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
            value_range=(1, MAXSPEED),
            click_increment=1,
            manager=manager,
            container=uiscreen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        # Void space
        self.elements["void_button"] = UIButton(
            relative_rect=pygame.Rect((x, topmargin+DIM["BUTTON_HEIGHT"]*5.6), (DIM["BUTTON_WIDTH"]*factor, DIM["BUTTON_HEIGHT"])),
            text="-",
            manager=manager,
            container=uiscreen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        self.elements["console"] = ui_text_box.UITextBox(
            relative_rect=pygame.Rect((0, 0), (DIM["CONSOLE_WIDTH"], DIM["SCREEN_HEIGHT"]-DIM["GAME_SCREEN_HEIGHT"])),
            html_text="",
            container=consoleContainer,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            },
            object_id="#console"
        )
        self.elements["stats"] = ui_text_box.UITextBox(
            relative_rect=pygame.Rect((DIM["CONSOLE_WIDTH"], 0), (DIM["GAME_SCREEN_WIDTH"]-DIM["CONSOLE_WIDTH"], DIM["SCREEN_HEIGHT"]-DIM["GAME_SCREEN_HEIGHT"])),
            html_text="",
            container=consoleContainer,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            },
            object_id="#stats"
        )
        self.stats = {
            "Active Bikes" : 0,
            "Active Lootboxes" : 0,
            "Alive Agents" : 0,
            "Dead Agents" : 0,
            "Void Agents" : 0
        }
        self.agents = {}
        return self.elements

    def render_game_screen(self, screen:UIContainer) -> None:
        """
        Render the game screen
        """
        screen.fill((255, 255, 255))
        self.draw_grid(screen)
        # Draw awdi
        self.awdi.draw(screen, self.offsetX, self.offsetY, self.zoom)
        # # Draw lootboxes
        for lootbox in self.lootboxes.values():
            lootbox.draw(screen, self.offsetX, self.offsetY, self.zoom)
        # Draw agents
        for bike in self.bikes.values():
            bike.draw(screen, self.offsetX, self.offsetY, self.zoom)
        # Draw overlays
        for bike in self.bikes.values():
            bike.draw_overlay(screen, self.offsetX, self.offsetY, self.zoom)
        for lootbox in self.lootboxes.values():
            lootbox.draw_overlay(screen, self.offsetX, self.offsetY, self.zoom)
        self.awdi.draw_overlay(screen, self.offsetX, self.offsetY, self.zoom)
        screen.blit(self.mouseSurface, (OVERLAY["FPS_PAD"], OVERLAY["FPS_PAD"]))
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
                        if event.pos[0] < DIM["GAME_SCREEN_WIDTH"] and event.pos[1] < DIM["GAME_SCREEN_HEIGHT"]:
                            self.dragging = True
                            self.mouseX, self.mouseY = event.pos
                            self.propagate_click(event.pos)
                    case 4:  # Scroll up
                        if event.pos[0] < DIM["GAME_SCREEN_WIDTH"] and event.pos[1] < DIM["GAME_SCREEN_HEIGHT"]:
                            self.adjust_zoom(1.1, event.pos)
                    case 5:  # Scroll down
                        if event.pos[0] < DIM["GAME_SCREEN_WIDTH"] and event.pos[1] < DIM["GAME_SCREEN_HEIGHT"]:
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
            # Handle auto play
            case self.playEvent:
                self.elements["round_slider"].set_current_value(self.round + 1)
                self.change_round(self.round + 1)
                if self.round == self.maxRound:
                    self.toggle_play()
            #Draw mouse coordinates
            case self.mouseDraw:
                self.draw_mouse_coords()

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
        self.playSpeed = min(MAXSPEED, max(1, newSpeed))
        self.elements["play_pause_speed_slider"].set_current_value(self.playSpeed)
        self.elements["play_pause_speed"].set_text(f"{self.playSpeed} Iteration/Sec")
        if self.isPlaying:
            pygame.time.set_timer(self.playEvent, int(1000//self.playSpeed))

    def change_round(self, newRound:int) -> None:
        """
        Change the current round using the round controls and json data
        Determine statistics
        """
        self.round = max(0, min(self.maxRound, newRound))
        self.elements["round_count"].set_text(f"Iteration: {self.round}")
        self.elements["iteration_count"].set_text(f"Round: {self.round // ITERATIONLENGTH}")
        self.elements["console"].html_text = ""
        self.stats["Void Agents"] = 0
        #Reload bikes
        self.update_bikes()
        self.update_lootboxes()
        if self.awdi is None:
            self.awdi = Awdi(self.jsonData[self.round]["audi"], self.jsonData[self.round]["bikes"])
        else:
            self.awdi.update_awdi(self.jsonData[self.round]["audi"], self.jsonData[self.round]["bikes"])
        if self.round == 0:
            self.aliveAgents = len(self.agents.values())
            self.stats["Dead Agents"] = 0
        self.update_stats()
        self.elements["stats"].rebuild()
        self.elements["console"].rebuild()

    def update_lootboxes(self) -> None:
        """
        Update the lootboxes
        """
        # Reload lootboxes
        lootboxes = {}
        for lootboxid, lootboxdata in self.jsonData[self.round]["loot_boxes"].items():
            if lootboxid not in self.lootboxes:
                lootboxObj = Lootbox(lootboxid, lootboxdata)
            else:
                lootboxObj = self.lootboxes[lootboxid]
                lootboxObj.update_lootbox(lootboxdata)
            lootboxes[lootboxid] = lootboxObj
        self.compare_lootboxes(lootboxes)

    def update_bikes(self) -> None:
        """
        Update the bikes
        """
        bikes = {}
        bikeAgents = {}
        for bikeid, bikedata in self.jsonData[self.round]["bikes"].items():
            #Determine next position for arrow drawing
            nextRound = self.jsonData[min(self.round+1, self.maxRound)]["bikes"]
            if bikeid in nextRound:
                nextPos = nextRound[bikeid]["physical_state"]["position"]["x"], nextRound[bikeid]["physical_state"]["position"]["y"]
                if "orientation" in nextRound[bikeid]:
                    nextOrient = nextRound[bikeid]["orientation"]
                else:
                    nextOrient = -2
            else:
            #If bike is dead, draw arrow to empty
                nextPos = 0, 0
                nextOrient = 0
            # If bike is new, create new bike
            if bikeid not in self.bikes:
                bikeObj = Bike(bikeid, bikedata, self.jsonData[self.round]["agents"], nextPos, nextOrient, self.jsonData[min(self.round+1, self.maxRound)]["agents"])
            else:
                bikeObj = self.bikes[bikeid]
                bikeObj.update_bike(bikedata, self.jsonData[self.round]["agents"], nextPos, nextOrient, self.jsonData[min(self.round+1, self.maxRound)]["agents"])
            bikes[bikeid] = bikeObj
            bikeAgents.update(bikeObj.get_agents())
        self.compare_bikes(bikes)
        self.compare_agents(self.jsonData[self.round]["agents"], bikeAgents)

    def propagate_click(self, mousePos:tuple) -> None:
        """
        Propagate the click to all entities
        """
        mouseX, mouseY = mousePos
        for bike in self.bikes.values():
            bike.propagate_click(mouseX, mouseY, self.zoom)
        for lootbox in self.lootboxes.values():
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
        height = DIM["GAME_SCREEN_HEIGHT"]
        # Draw vertical lines
        for x in range(-int(zoomedSpacing) + int(startX), width, int(zoomedSpacing)):
            pygame.draw.line(surface, BGCOLOURS["GRID"], (x, 0), (x, height))
        # Draw horizontal lines
        for y in range(-int(zoomedSpacing) + int(startY), height, int(zoomedSpacing)):
            pygame.draw.line(surface, BGCOLOURS["GRID"], (0, y), (width, y))

    def draw_mouse_coords(self) -> None:
        """
        Draw the mouse coordinates on the game screen
        """
        font = pygame.font.SysFont("Arial", 15)
        x = round((self.mouseXCur / self.zoom - self.offsetX / self.zoom) / COORDINATESCALE, PRECISION)
        y = round((self.mouseYCur / self.zoom - self.offsetY / self.zoom) / COORDINATESCALE, PRECISION)
        self.mouseSurface = font.render(f"({x}, {y})", True, (0, 0, 0))

    def log(self, message:str, tag=None) -> None:
        """
        Draw the console on the game screen
        """
        if tag is None:
            tag = "DEFAULT"
        color = CONSOLE[tag]
        # Format message as HTML with reduced margin in paragraph tag
        text = f'<p style="margin: 0; line-height: 100%;"><font color="{color}">{message}</font></p>'
        newText = self.elements["console"].html_text + text
        self.elements["console"].html_text = newText

    def compare_agents(self, newAgents:dict, bikeAgents:dict) -> None:
        """
        Compare the agents to the previous round and update console and stats
        """
        dead = 0
        for agentid, agent in self.agents.items():
            if agentid not in newAgents:
                # Agent has died from exhaustion
                groupID = agent["group_id"]
                if agent["energy_level"] < ENERGYTHRESHOLD:
                    self.log(f"Agent {agentid} ({groupID}) has run out of energy!", "ERROR")
                # Agent has died from being run over
                elif (pow(agent["location"]["x"]-self.awdi.x, 2) < pow(EPSILON, 2)) and (pow(agent["location"]["y"]-self.awdi.y, 2) < pow(EPSILON, 2)):
                    self.log(f"Agent {agentid} ({groupID}) has been run over by the Owdi!", "ERROR")
                # Agent has died for unknown reasons
                else:
                    self.log(f"Agent {agentid} ({groupID}) has died for unknown reasons!", "ERROR")
                dead += 1
            else:
                #Check if agent has moved bikes - on bike parameter is not reliable
                if agentid not in bikeAgents:
                    self.log(f"Agent {agentid} is in the void.", "INFO")
                    self.stats["Void Agents"] += 1
        self.stats["Alive Agents"] = len(newAgents.values())
        self.stats["Dead Agents"] = self.aliveAgents - self.stats["Alive Agents"]
        self.agents = newAgents

    def compare_bikes(self, newBikes:dict) -> None:
        """
        Compare the agents to the previous round and update console and stats
        """
        activeBikes = 0
        for bikeid, _ in self.bikes.items():
            if bikeid not in newBikes:
                self.log(f"Bike {bikeid} has died!", "ERROR")
            elif len(newBikes[bikeid].get_agents()) > 0:
                activeBikes += 1
        self.stats["Active Bikes"] = activeBikes
        self.bikes = newBikes

    def compare_lootboxes(self, newLootboxes:dict) -> None:
        """
        Compare the lootboxes to the previous round and update console and stats
        """
        for lootboxid, _ in self.lootboxes.items():
            if lootboxid not in newLootboxes:
                self.log(f"Lootbox {lootboxid} has been collected!", "INFO")
        self.stats["Active Lootboxes"] = len(newLootboxes.values())
        self.lootboxes = newLootboxes

    def update_stats(self) -> None:
        """
        Update the stats
        """
        text = ""
        for stat, value in self.stats.items():
            text += f"<font face=verdana size=3 color=#FFFFFF><b>{stat}</b>: {value}</font><br>"
        self.elements["stats"].html_text = text

    def set_json(self, data:dict) -> None:
        """
        Set the data for the game screen
        """
        if data is None:
            return
        self.jsonData = data
        self.maxRound = len(data) - 1
