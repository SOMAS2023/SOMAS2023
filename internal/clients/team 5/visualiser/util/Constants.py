"""
Constants used for the visualiser.
"""
# Screen constants
WINDOW_TITLE = "SOMAS Visualiser"
FRAMERATE = 60
MINZOOM, MAXZOOM, ZOOM = 0.2, 2.5, 0.5
COORDINATESCALE = 1.5
OVERLAY = {
    "FONT" : "Arial",
    "WIDTH": 150,
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
}
BIKE = {
    "LINE_WIDTH": 1,
    "LINE_COLOUR": "#000000",
    "COLOURS": {
        "MINHUE" : 240,
        "MAXHUE" : 300,
        "MINSAT" : 50,
        "MAXSAT" : 80,
        "MINVAL" : 60,
        "MAXVAL" : 80,
    },
}
AWDI = {
    "COLOUR" : "#B0B0B0",
    "LINE_WIDTH": 2,
    "LINE_COLOUR": "#000000",
    "FONT_SIZE": 20,
    "SIZE": 100,
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
    "Red": (255, 0, 0),
    "Green": (0, 128, 0),
    "Blue": (0, 0, 255),
    "Yellow": (255, 255, 0),
    "Orange": (255, 165, 0),
    "Purple": (128, 0, 128),
    "Pink": (255, 192, 203),
    "Brown": (165, 42, 42),
    "Gray": (128, 128, 128),
    "White": (255, 255, 255)
}
