package levels

// LevelManager handles level switching
type LevelManager struct {
	ActiveLevel *Level
	LevelIDs    []LevelID // all levels available, that can be loaded
}

func NewLevelManager() *LevelManager {
	return &LevelManager{
		ActiveLevel: nil,
	}
}

// // scan all levels in the level_files folder, and
// // add them to the LevelIDs slice
// func (l *levelManager) ScanAvailableLevels() {

// }

func (l *LevelManager) LoadLevelByID(id LevelID) bool {
	if l.ActiveLevel != nil {
		l.ActiveLevel.Unload()
	}
	l.ActiveLevel = NewLevel(id)
	return l.ActiveLevel.Load(nil) // syncronous level load (TODO: make it async)
}
