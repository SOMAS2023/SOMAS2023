import tkinter as tk
from tkinter import filedialog
import pygame
import pygame_gui
from pygame_gui.elements import UIButton
from visualiser.util.Constants import DIM
from visualiser.util.HelperFunc import make_center

class MainMenuScreen:
    def init_ui(self, elements:dict, manager:pygame_gui.UIManager, screen:pygame_gui.core.UIContainer) -> dict:
        """
        Initialise the UI for the main menu screen
        """
        x, y = make_center((DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"]), (DIM["UI_WIDTH"], DIM["SCREEN_HEIGHT"]))
        elements["load_json"] = UIButton(
            relative_rect=pygame.Rect((x, 2*y-40), (DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"])),
            text="Load JSON",
            manager=manager,
            container=screen,
            anchors={
                "left": "left",
                "right": "left",
                "top": "top",
                "bottom": "top",
            }
        )
        return elements

    def handle_events(self, elements:dict, event:pygame.event.Event, switch) -> None:
        """
        Handle events for the main menu screen
        """
        if event.ui_element == elements["load_json"]:
            switch("game_screen")