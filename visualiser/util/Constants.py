"""
Constants used for the visualiser.
"""
# Screen constants
WINDOW_TITLE = "SOMAS Visualiser"
FRAMERATE = 60
MINZOOM, MAXZOOM, ZOOM = 0.2, 2.5, 0.5
COORDINATESCALE = 15
PRECISION = 2
# JSONPATH = "./internal/clients/team_5/visualiser/json/test2.json"
JSONPATH = "game_dump.json"
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
    "TRANSPARENCY": 240,
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
    "UI_WIDTH": 280,
    "BUTTON_WIDTH": 220,
    "BUTTON_HEIGHT": 70,
    "GRIDSPACING": 20,
}
THEMEJSON = "visualiser/theme.json"
BGCOLOURS = {
    "GUI" : "#E0E0E0",
    "MAIN" : "#F0F0F0",
    "GRID" : "#E5E5E5",
}
COLOURS = {
    "red": "#E05558",
    "orange": "#D5C801",
    "yellow": "#D5C801",
    "green": "#7BBD01",
    "blue": "#5E82FD",
    "purple": "#A575ED",
    "pink": "#DE82C3",
    "brown": "#AC6223",
    "gray": "#666666",
    "white": "#FFFFFF"
}
