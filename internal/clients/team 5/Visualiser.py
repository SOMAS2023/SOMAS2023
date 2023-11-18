"""
Visualiser for SOMAS world
"""
# pylint: disable=no-member
import pygame
from pygame_gui import UIManager
from pygame_gui.elements import UIButton, UIImage
from pygame_gui.core import UIContainer
import visualiser.Constants as CONST
from visualiser.entities import Agents, Lootboxes, Awdi, Bikes
from visualiser.HelperFunc import pick_random_colour

class Visualiser:
    def __init__(self) -> None:
        self.screenState = "main_menu"
        self.UIState = "main_menu"
        self.running = True
        pygame.init()
        # Set screens, UI manager, caption, clock and timeDelta
        self.window = pygame.display.set_mode((CONST.SCREEN_WIDTH, CONST.SCREEN_HEIGHT))
        self.manager = UIManager((CONST.SCREEN_WIDTH, CONST.SCREEN_HEIGHT))
        self.gamescreen = pygame.Surface((CONST.GAME_SCREEN_WIDTH, CONST.SCREEN_HEIGHT))
        self.UIscreen = UIContainer(relative_rect=pygame.Rect((CONST.GAME_SCREEN_WIDTH, 0), 
                          (CONST.SCREEN_WIDTH - CONST.GAME_SCREEN_WIDTH, CONST.SCREEN_HEIGHT)), 
                          manager=self.manager)
        pygame.display.set_caption(CONST.WINDOW_TITLE)
        self.clock = pygame.time.Clock()
        # UI element dictionaries
        self.UIElements = {
            "main_menu": {},
            "game_screen": {},
            "dead_players": {}
        }
        # Initialise UI elements
        self.init_ui()

    def init_ui(self) -> None:
        """
        Initialise UI elements
        """
        # Common UI elements
        self.UIbackground = pygame.Surface((CONST.SCREEN_WIDTH - CONST.GAME_SCREEN_WIDTH, CONST.SCREEN_HEIGHT))
        self.UIbackground.fill(CONST.BGCOLOURS["GUI"])
        # Background colour
        UIImage(relative_rect=pygame.Rect((0, 0), (CONST.SCREEN_WIDTH - CONST.GAME_SCREEN_WIDTH, CONST.SCREEN_HEIGHT)),
                image_surface=self.UIbackground,
                manager=self.manager,
                container=self.UIscreen)
        # Main menu UI elements
        # Load JSON button
        self.UIElements["main_menu"]["load_json"] = UIButton(
            relative_rect=pygame.Rect((50, 50), (200, 50)),
            text="Load JSON",
            manager=self.manager,
            container=self.UIscreen
        )

    def run_loop(self) -> None:
        """
        Main loop for visualiser
        """
        while self.running:
            timeDelta = self.clock.tick(CONST.FRAMERATE) / 1000.0
            self.handle_events()
            # Update the game screen
            match self.screenState:
                case "main_menu":
                    self.main_menu_screen()
                case _:
                    self.main_menu_screen()
            # Update the UI panel
            match self.UIState:
                case "main_menu":
                    self.main_menu_control_panel()
                case _:
                    self.main_menu_control_panel()
            # Update the display
            self.manager.update(timeDelta)
            self.window.blit(self.gamescreen, (0, 0))
            self.manager.draw_ui(self.window)
            pygame.display.flip()
        pygame.quit()

    def handle_events(self) -> None:
        """
        Handle events in the visualiser
        """
        for event in pygame.event.get():
            match event.type:
                case pygame.QUIT:
                    self.running = False
                case _:
                    pass
            self.manager.process_events(event)

    def main_menu_screen(self) -> None:
        """
        UI logic for the main menu screen
        """
        self.gamescreen.fill(pick_random_colour(CONST.COLOURS))



    def main_menu_control_panel(self) -> None:
        # self.UIscreen.fill(CONST.BGCOLOURS["GUI"])
        pass

    def game_screen(self) -> None:
        pass

    def dead_player_screen(self) -> None:
        pass

    def game_control_panel(self) -> None:
        pass


if __name__ == "__main__":
    visualiser = Visualiser()
    visualiser.run_loop()
