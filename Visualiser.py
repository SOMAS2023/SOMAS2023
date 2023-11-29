"""
Visualiser for SOMAS world
"""
# pylint: disable=no-member, import-error, no-name-in-module
import tkinter as tk
from tkinter import filedialog
import json
from os.path import exists
import pygame
import pygame_gui
from pygame_gui import UIManager
from pygame_gui.elements import UIButton, UIImage
from pygame_gui.core import UIContainer
from visualiser.util.Constants import WINDOW_TITLE, FRAMERATE, DIM, BGCOLOURS, THEMEJSON, OVERLAY, JSONPATH
from visualiser.util.HelperFunc import make_center
from visualiser.GameScreen import GameScreen

class Visualiser:
    def __init__(self) -> None:
        self.screenState = "main_menu"
        self.UIState = "main_menu"
        self.running = True
        self.jsondata = None
        pygame.init()
        # Set screens, UI manager, caption, clock and timeDelta
        self.window = pygame.display.set_mode((DIM["SCREEN_WIDTH"], DIM["SCREEN_HEIGHT"]))
        self.manager = UIManager((DIM["SCREEN_WIDTH"], DIM["SCREEN_HEIGHT"]), THEMEJSON)
        self.gamescreen = pygame.Surface((DIM["GAME_SCREEN_WIDTH"], DIM["SCREEN_HEIGHT"]))
        self.UIscreen = UIContainer(relative_rect=pygame.Rect((DIM["GAME_SCREEN_WIDTH"], 0),
                          (DIM["UI_WIDTH"], DIM["SCREEN_HEIGHT"],)),
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
            self.handle_events()
            # Update the display
            match self.screenState:
                case "main_menu":
                    self.render_main_menu()
                case "game_screen":
                    self.gameScreenManager.render_game_screen(self.gamescreen)
            self.manager.update(timeDelta)
            self.window.blit(self.gamescreen, (0, 0))
            self.manager.draw_ui(self.window)
            self.draw_fps()
            pygame.display.flip()
        pygame.quit()

    def draw_fps(self) -> None:
        """
        Draw the FPS counter
        """
        fps = self.clock.get_fps()
        font = pygame.font.SysFont("Arial Narrow", 20)
        surface = font.render(f"FPS: {fps:.2f}", True, "#555555")
        self.window.blit(surface, (DIM["GAME_SCREEN_WIDTH"]-surface.get_width()-OVERLAY["FPS_PAD"], OVERLAY["FPS_PAD"]))

    def handle_events(self) -> None:
        """
        Handle events in the visualiser
        """
        for event in pygame.event.get():
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
            # Handle UI events
            match self.screenState:
                case "main_menu":
                    self.process_main_menu_events(event)
                case "game_screen":
                    self.process_game_screen_events(event)
            self.manager.process_events(event)

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
                    self.json_parser(filepath)
                    self.gameScreenManager.change_round(0)
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
        self.UIElements["game_screen"] = self.gameScreenManager.init_ui(self.manager, self.UIscreen)
    def test(self) -> None:
        """
        Test function
        """
        if exists(JSONPATH):
            self.json_parser(JSONPATH)
            self.gameScreenManager.change_round(0)
            self.run_loop("game_screen")
        else:
            self.run_loop()

if __name__ == "__main__":
    visualiser = Visualiser()
    visualiser.test()
    # visualiser.run_loop()
