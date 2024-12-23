package levels

var levelMgrSingleton *levelManager = nil

type levelManager struct {
	ActiveLevel *Level
	LevelIDs    []LevelID // all levels available, that can be loaded
}

// returns a reference to the level manager singleton instance. it gets initialized if it isn't already
func Manager() *levelManager {
	if levelMgrSingleton == nil {
		levelMgrSingleton = &levelManager{
			ActiveLevel: nil,
		}
	}
	return levelMgrSingleton
}

// // scan all levels in the level_files folder, and
// // add them to the LevelIDs slice
// func (l *levelManager) ScanAvailableLevels() {

// }

func (l *levelManager) LoadLevelByID(id LevelID) {
	if l.ActiveLevel != nil {
		l.ActiveLevel.Unload()
	}
	l.ActiveLevel = NewLevel(id)
	l.ActiveLevel.Load(nil) // syncronous level load (TODO: make it async)

}
