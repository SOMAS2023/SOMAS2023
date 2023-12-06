"""
Constants used for the visualiser.
"""
# Screen constants
WINDOW_TITLE = "SOMAS Visualiser"
FRAMERATE = 60
MINZOOM, MAXZOOM, ZOOM = 0.2, 2.5, 0.5
COORDINATESCALE = 15
PRECISION = 2
EPSILON = 8
ENERGYTHRESHOLD = 0.1
THEMEJSON = "visualiser/theme.json"
JSONPATH = "game_dump.json"
MAXSPEED = 50
ITERATIONLENGTH = 100
TEXT = {
    "FONT" : "Arial",
    "FONT_SIZE": 1,
    "PADDING": 2,
    "DEFAULT_COLOUR" : "#FFFFFF",
    "LINE_SPACING": 2,
    "BACKGROUND_COLOUR": "#699ff5",
    "TEXT_COLOUR": "#FFFFFF",
    "LINE_COLOUR": "#d8e0ed",
    "LINE_WIDTH": 2,
    "BORDER_WIDTH": 2,
    "BORDER_COLOUR": "#000000",
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
    "COLOURS": {
        "MINHUE" : 300,
        "MAXHUE" : 300,
        "MINSAT" : 0,
        "MAXSAT" : 0,
        "MINVAL" : 60,
        "MAXVAL" : 80,
    },
    "TRANSPARENCY": 150,
}
AWDI = {
    "COLOUR" : "#0F0F0F",
    "LINE_WIDTH": 2,
    "LINE_COLOUR": "#000000",
    "FONT_SIZE": 30,
    "SIZE": 60,
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
    "HEIGHT" : 30,
    "WIDTH": 120,
    "LINE_WIDTH": 2,
    "LINE_COLOUR": "#000000",
    "FONT_SIZE": 20,
}
DIM = {
    "SCREEN_WIDTH": 1280,
    "SCREEN_HEIGHT": 720,
    "GAME_SCREEN_WIDTH": 1000,
    "GAME_SCREEN_HEIGHT": 580,
    "UI_WIDTH": 280,
    "BUTTON_WIDTH": 220,
    "BUTTON_HEIGHT": 70,
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
