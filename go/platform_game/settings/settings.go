package settings

const (
	MAP_ROWS   = 11
	TITE_SIZE  = 64
	MAP_HEIGHT = MAP_ROWS * TITE_SIZE
)

var levelMap = [MAP_ROWS]string{
	"                                  ",
	"                                  ",
	"                                  ",
	" P                                ",
	"TT                                ",
	"                                  ",
	"                                  ",
	"                                  ",
	"           T     TTTT  T          ",
	"          TT T            T       ",
	" TTTTTTTTTTT TTTTTTTTTT  TTTTTTT T",
}

func GetLevelMap() []string {
	level := make([]string, 0, len(levelMap))
	for _, row := range levelMap {
		level = append(level, row)
	}

	return level
}

const (
	SCREEN_WIDTH  = 1200
	SCREEN_HEIGHT = MAP_HEIGHT
	WINDOW_TITLE  = "The Fucking World"
)
