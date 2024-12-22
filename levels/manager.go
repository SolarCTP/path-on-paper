package levels

const (
	levelFilesPath string = "./level_files"
)

type LevelManager struct {
	ActiveLevel *Level
	LevelIDs    []LevelID // all levels available, that can be loaded
}

func NewLevelManager() *LevelManager {
	return &LevelManager{
		ActiveLevel: nil,
	}
}

// scan all levels in the level_files folder, and
// add them to the LevelIDs slice
func (l *LevelManager) ScanAvailableLevels() {

}

func (l *LevelManager) LoadLevelByID(id LevelID) {
	l.ActiveLevel.unload()

}
