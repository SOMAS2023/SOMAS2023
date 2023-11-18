import pygame
import pygame_gui
from pygame_gui.elements import UIButton, UILabel
from visualiser.util.Constants import DIM
from visualiser.util.HelperFunc import make_center

class GameScreen:
    def init_ui(self, elements:dict, manager:pygame_gui.UIManager, screen:pygame_gui.core.UIContainer) -> dict:
        """
        Initialise the UI for the main menu screen
        """
        x, _ = make_center((DIM["BUTTON_WIDTH"], DIM["BUTTON_HEIGHT"]), (DIM["UI_WIDTH"], DIM["SCREEN_HEIGHT"]))
        elements["reset"] = UIButton(
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
        elements["round_count"] = UILabel(
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
        elements["increase_round"] = UIButton(
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
        elements["decrease_round"] = UIButton(
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
        return elements
