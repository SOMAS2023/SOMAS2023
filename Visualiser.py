"""
Visualiser for SOMAS world
"""
# pylint: disable=no-member, import-error, no-name-in-module, pointless-string-statement
import tkinter as tk
from tkinter import filedialog
import json
from os.path import exists
import sys
import pygame
import pygame_gui
from pygame_gui import UIManager
from pygame_gui.elements import UIButton, UIImage
from pygame_gui.core import UIContainer
from visualiser.util.Constants import WINDOW_TITLE, FRAMERATE, DIM, BGCOLOURS, THEMEJSON, OVERLAY, JSONPATH, FPSDISPLAYRATE
from visualiser.util.HelperFunc import make_center
from visualiser.GameScreen import GameScreen

class Visualiser:
    def __init__(self) -> None:
        self.screenState = "main_menu"
        self.UIState = "main_menu"
        self.running = True
        self.jsondata = None
        self.fps = None
        pygame.init()
        self.drawInfo = pygame.USEREVENT + 101
        pygame.time.set_timer(self.drawInfo, 1000 // FPSDISPLAYRATE)
        # Set screens, UI manager, caption, clock and timeDelta
        self.window = pygame.display.set_mode((DIM["SCREEN_WIDTH"], DIM["SCREEN_HEIGHT"]))
        self.manager = UIManager((DIM["SCREEN_WIDTH"], DIM["SCREEN_HEIGHT"]), sys.argv[0] + "/../"+THEMEJSON)
        self.gamescreen = pygame.Surface((DIM["GAME_SCREEN_WIDTH"], DIM["SCREEN_HEIGHT"]))
        self.UIscreen = UIContainer(relative_rect=pygame.Rect((DIM["GAME_SCREEN_WIDTH"], 0),
                        (DIM["UI_WIDTH"], DIM["SCREEN_HEIGHT"],)),
                        manager=self.manager)
        self.consoleContainer = UIContainer(relative_rect=pygame.Rect((0, DIM["GAME_SCREEN_HEIGHT"]),
                        (DIM["GAME_SCREEN_WIDTH"], DIM["SCREEN_HEIGHT"]-DIM["GAME_SCREEN_HEIGHT"])),
                        manager=self.manager)
        pygame.display.set_caption(WINDOW_TITLE)
        self.clock = pygame.time.Clock()
        # UI element dictionaries
        self.UIElements = {
            "main_menu": {},
            "game_screen": {},
            "dead_players": {}
        }
        # Initialise UI elements
        self.gameScreenManager = GameScreen()
        self.init_ui()

    def init_ui(self) -> None:
        """
        Initialise UI elements.
        """
        # Common UI elements
        self.UIbackground = pygame.Surface((DIM["UI_WIDTH"], DIM["SCREEN_HEIGHT"]))
        self.UIbackground.fill(BGCOLOURS["GUI"])
        # Background colour
        UIImage(relative_rect=pygame.Rect((0, 0), (DIM["UI_WIDTH"], DIM["SCREEN_HEIGHT"])),
                image_surface=self.UIbackground,
                manager=self.manager,
                container=self.UIscreen)
        # Main menu UI elements
        buttonWidth, buttonHeight = 220, 70
        x, y = make_center((buttonWidth, buttonHeight), (DIM["UI_WIDTH"], DIM["SCREEN_HEIGHT"]))
        # Load JSON button
        self.UIElements["main_menu"]["load_json"] = UIButton(
            relative_rect=pygame.Rect((x, 2*y-40), (DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"])),
            text="Load JSON",
            manager=self.manager,
            container=self.UIscreen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        self.draw_fps()

    def switch_screen(self, newScreen:str) -> None:
        """
        Switch between screens, hide all elements in the current screen and show all elements
        in the new screen.
        """
        for screen, elements in self.UIElements.items():
            for _, obj in elements.items():
                if screen == newScreen:
                    obj.show()
                else:
                    obj.hide()
        self.screenState = newScreen

    def run_loop(self, screen:str="main_menu") -> None:
        """
        Main loop for visualiser.
        """
        self.switch_screen(screen)
        while self.running:
            timeDelta = self.clock.tick(FRAMERATE) / 1000.0
            # Update the display
            self.handle_events()
            match self.screenState:
                case "main_menu":
                    self.render_main_menu()
                case "game_screen":
                    self.gameScreenManager.render_game_screen(self.gamescreen)
            self.manager.update(timeDelta)
            self.window.blit(self.gamescreen, (0, 0))
            self.manager.draw_ui(self.window)
            # Draw FPS counter
            self.window.blit(self.fps, (DIM["GAME_SCREEN_WIDTH"]-self.fps.get_width()-OVERLAY["FPS_PAD"], OVERLAY["FPS_PAD"]))
            pygame.display.flip()
        pygame.quit()

    def draw_fps(self) -> None:
        """
        Draw the FPS counter
        """
        fps = self.clock.get_fps()
        font = pygame.font.SysFont("Arial Narrow", 20)
        self.fps = font.render(f"FPS: {fps:.2f}", True, "#555555")

    def handle_events(self) -> None:
        """
        Handle events in the visualiser
        """
        for event in pygame.event.get():
            # Handle UI events
            match self.screenState:
                case "main_menu":
                    self.process_main_menu_events(event)
                case "game_screen":
                    self.process_game_screen_events(event)
            self.manager.process_events(event)
            match event.type:
                # Quit the game
                case pygame.QUIT:
                    self.running = False
                # Handle key presses
                case pygame.KEYDOWN:
                    match event.key:
                        # Quit the game
                        case pygame.K_ESCAPE:
                            self.running = False
                case self.drawInfo:
                    self.draw_fps()

    def render_main_menu(self) -> None:
        """
        UI logic for the main menu screen
        """
        self.gamescreen.fill(BGCOLOURS["MAIN"])
        font = pygame.font.SysFont("Arial Narrow", 90)
        surface = font.render("SOMAS Visualiser", True, "#555555")
        textWidth, textHeight = surface.get_size()
        x, y = make_center((textWidth, textHeight), (DIM["GAME_SCREEN_WIDTH"], DIM["SCREEN_HEIGHT"]))
        pygame.draw.line(self.gamescreen, "#555555", (x, y+textHeight), (x+textWidth, y+textHeight), 2)
        # Divider line at edge of UI screen
        lineWidth = 1
        pygame.draw.line(self.gamescreen, "#555555", (DIM["GAME_SCREEN_WIDTH"]-lineWidth, 0), (DIM["GAME_SCREEN_WIDTH"]-lineWidth, DIM["SCREEN_HEIGHT"]), lineWidth)
        self.gamescreen.blit(surface, (x, y))

    def process_main_menu_events(self, event:pygame.event.Event) -> None:
        """
        Process events in the main menu screen
        """
        elements = self.UIElements["main_menu"]
        if event.type == pygame_gui.UI_BUTTON_PRESSED:
            if event.ui_element == elements["load_json"]:
                #Load JSON for game using tkinter
                root = tk.Tk()
                root.withdraw()
                filepath = filedialog.askopenfilename(
                    initialdir=JSONPATH,
                    title="Select JSON file",
                    filetypes=(("JSON files", "*.json"), ("all files", "*.*"))
                )
                root.destroy()
                if filepath != "":
                    try:
                        self.load_game(filepath)
                    except: # pylint: disable=bare-except
                        print("Attempted to load incompatible/outdated JSON file.")
                        self.switch_screen("main_menu")

    def load_game(self, filepath:str) -> None:
        """
        Load the game from the JSON file
        """
        self.json_parser(filepath)
        self.gameScreenManager.change_iteration(0)
        self.gameScreenManager.log("Welcome to the visualiser!")
        self.gameScreenManager.log(f"Max Iterations: {(self.gameScreenManager.maxRound * self.gameScreenManager.roundLength)-1}", "INFO")
        self.gameScreenManager.log(f"Max Rounds: {self.gameScreenManager.maxRound}", "INFO")
        self.gameScreenManager.log(f"There are {self.gameScreenManager.roundLength} iterations per round.", "INFO")
        self.gameScreenManager.elements["console"].rebuild()
        self.switch_screen("game_screen")

    def process_game_screen_events(self, event:pygame.event.Event) -> None:
        """
        Process events in the main menu screen
        """
        elements = self.UIElements["game_screen"]
        if event.type == pygame_gui.UI_BUTTON_PRESSED:
            if event.ui_element == elements["reset"]:
                self.switch_screen("main_menu")
                self.jsondata = None
        self.gameScreenManager.process_events(event)

    def json_parser(self, filepath:str) -> None:
        """
        Reads the simulated JSON file and stores the data
        """
        with open(filepath, "r", encoding="utf-8") as f:
            data = json.load(f)
        self.jsondata = data
        self.gameScreenManager.set_json(data)
        self.UIElements["game_screen"] = self.gameScreenManager.init_ui(self.manager, self.UIscreen, self.consoleContainer)

    def start(self) -> None:
        """
        Start function
        """
        filepath = sys.argv[0] + "/../" + JSONPATH
        if exists(filepath):
            # Load game from JSON file
            try:
                self.load_game(filepath)
                self.run_loop("game_screen")
            except: # pylint: disable=bare-except
                print("Attempted to load incompatible/outdated JSON file.")
                self.run_loop()
        else:
            self.run_loop()

if __name__ == "__main__":
    visualiser = Visualiser()
    # Run profiler to check for optimisations
    OPTIM = False
    if OPTIM:
        import cProfile
        import subprocess
        profiler = cProfile.Profile()
        profiler.enable()
        profiler.run("visualiser.start()")
        profiler.dump_stats('visualiser/profiles/stats.prof')
        subprocess.Popen("snakeviz visualiser/profiles/stats.prof", shell=True)
    else:
        visualiser.start()

"""
TODO:
-Motivation
-Design decisions
    - Why i selected certain attributes
        - 
"""