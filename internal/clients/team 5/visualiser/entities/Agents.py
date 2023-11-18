import pygame

class Agent(pygame.sprite.Sprite):
    def __init__(self, colour) -> None:
        super(Agent, self).__init__()
        self.boxSize = 20
        self.colour = colour
        self.surface = pygame.Surface((self.boxSize, self.boxSize))
        self.box = self.surface.get_rect()
