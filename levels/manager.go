package levels

import (
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

// LevelManager handles level switching
type LevelManager struct {
	ActiveLevel *Level
	LevelIDs    []LevelID // all levels available, that can be loaded
}

func NewLevelManager() *LevelManager {

	return &LevelManager{
		ActiveLevel: nil,
		LevelIDs:    searchAvailableLevels(),
	}
}

// searchAvailableLevels looks for valid levels inside of LEVELS_PATH
// and returns their IDs to
func searchAvailableLevels() []LevelID {
	var availableLevelIDs []LevelID

	dirEntries, err := os.ReadDir(LEVELS_PATH)
	if err != nil {
		log.Fatal("Could not find levels folder: " + LEVELS_PATH)
	}
	log.Println("Looking for levels in folder " + LEVELS_PATH)
	for _, entry := range dirEntries {
		levelIdStr, isPng := strings.CutSuffix(entry.Name(), ".png")
		if isPng && !entry.IsDir() {
			lvlID, err := strconv.ParseUint(levelIdStr, 10, 64)
			if err != nil {
				log.Println("Level", entry.Name(), "doesn't have a valid name. Can't read its ID.")
			} else {
				availableLevelIDs = append(availableLevelIDs, LevelID(lvlID))
				log.Println("Found level", entry.Name())
			}
		} else {
			log.Println(entry.Name(), "is not a level (not png or folder).")
		}
	}
	log.Println("Search done. Found", len(availableLevelIDs), "levels.")
	return availableLevelIDs
}

// // scan all levels in the level_files folder, and
// // add them to the LevelIDs slice
// func (l *levelManager) ScanAvailableLevels() {

// }

func (l *LevelManager) LoadLevelByID(id LevelID) bool {
	// check that the level exists
	if slices.Index(l.LevelIDs, id) == -1 {
		log.Println("Could not find level with ID", strconv.FormatUint(uint64(id), 10))
		return false
	}

	nextLevel := NewLevel(id)
	// TODO: load level concurrently
	if !nextLevel.Load(nil) {
		return false
	}
	if l.ActiveLevel != nil {
		l.ActiveLevel.Unload()
	}
	l.ActiveLevel = nextLevel
	return true
}

func (l *LevelManager) LoadNextAvailableLevel() bool {
	currentLvlIdIndex := slices.Index(l.LevelIDs, l.ActiveLevel.ID)
	var nextLevelIndex int = -1
	if currentLvlIdIndex == len(l.LevelIDs)-1 {
		nextLevelIndex = 0
	} else {
		nextLevelIndex = currentLvlIdIndex + 1
	}
	return l.LoadLevelByID(l.LevelIDs[nextLevelIndex])
}

func (l *LevelManager) LoadPrevAvailableLevel() bool {
	currentLvlIdIndex := slices.Index(l.LevelIDs, l.ActiveLevel.ID)
	var prevLevelIndex int = -1
	if currentLvlIdIndex == 0 {
		prevLevelIndex = len(l.LevelIDs) - 1
	} else {
		prevLevelIndex = currentLvlIdIndex - 1
	}
	return l.LoadLevelByID(l.LevelIDs[prevLevelIndex])
}
