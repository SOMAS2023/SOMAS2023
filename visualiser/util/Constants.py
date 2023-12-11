"""
Constants used for the visualiser.
"""
# Screen constants
WINDOW_TITLE = "SOMAS Visualiser"
FRAMERATE, FPSDISPLAYRATE = 75, 12
MINZOOM, MAXZOOM, ZOOM = 0.2, 2.5, 0.3
COORDINATESCALE = 20
PRECISION = 2
EPSILON = 12
ENERGYTHRESHOLD = 0.1
THEMEJSON = "visualiser/theme.json"
JSONPATH = "game_dump.json"
MAXSPEED = 100
ARROWS = {
    "NUM_ARROWS": 5,
    "ARROW_LENGTH": 15,
    "ARROW_ANGLE": 30,
}
CONSOLE = {
    "DEFAULT" : "#FFFFFF",
    "ERROR" : "#FF0000",
    "INFO" : "#FFFF00",
}
OVERLAY = {
    "FONT" : "Arial",
    "WIDTH": 300,
    "FONT_SIZE": 15,
    "PADDING": 2,
    "LINE_SPACING": 2,
    "BACKGROUND_COLOUR": "#699ff5",
    "TEXT_COLOUR": "#FFFFFF",
    "LINE_COLOUR": "#d8e0ed",
    "LINE_WIDTH": 2,
    "BORDER_WIDTH": 2,
    "BORDER_COLOUR": "#000000",
    "FPS_PAD" : 5,
}
BIKE = {
    "LINE_WIDTH": 1,
    "LINE_COLOUR": "#000000",
    "TRANSPARENCY": 150,
}
OWDI = {
    "COLOUR" : "#0F0F0F",
    "LINE_WIDTH": 2,
    "LINE_COLOUR": "#000000",
    "FONT_SIZE": 30,
    "SIZE": 140,
}
AGENT = {
    "SIZE": 10,
    "LINE_WIDTH": 2,
    "LINE_COLOUR": "#000000",
    "FONT_SIZE": 20,
    "PADDING": 2,
}
LOOTBOX = {
    "DEFAULT_COLOUR" : "#000000",
    "HEIGHT" : 60,
    "WIDTH": 120,
    "LINE_WIDTH": 2,
    "LINE_COLOUR": "#000000",
    "FONT_SIZE": 20,
}
DIM = {
    "SCREEN_WIDTH": 1280,
    "SCREEN_HEIGHT": 720,
    "GAME_SCREEN_WIDTH": 1000,
    "GAME_SCREEN_HEIGHT": 575,
    "UI_WIDTH": 280,
    "BUTTON_WIDTH": 220,
    "BUTTON_HEIGHT": 60,
    "GRIDSPACING": 20,
    "CONSOLE_WIDTH": 750,
}
BGCOLOURS = {
    "GUI" : "#E0E0E0",
    "MAIN" : "#F0F0F0",
    "GRID" : "#E5E5E5",
}
COLOURS = {
    "red": "#E05558",
    "orange": "#D57901",
    "yellow": "#D5C801",
    "green": "#7BBD01",
    "blue": "#5E82FD",
    "purple": "#A575ED",
    "pink": "#DE82C3",
    "brown": "#AC6223",
    "gray": "#666666",
    "white": "#FFFFFF"
}
GOVERNANCE = {
    0: ("Democracy", COLOURS["blue"]),
    1: ("Leadership", COLOURS["green"]),
    2: ("Dictatorship", COLOURS["red"]),
    3: ("Invalid", "#000000"),
}
