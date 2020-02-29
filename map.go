package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const mapURL = "https://raw.githubusercontent.com/rytrose/dist-midi/master/soundMap.json"

// SoundMapEntry is a description of the sound mapped to a MIDI note.
type SoundMapEntry struct {
	Title        string `json:"title"`
	HoldToPlay   bool   `json:"holdToPlay"`
	AllowPausing bool   `json:"allowPausing"`
	Loop         bool   `json:"loop"`
}

// SoundMap is a map from MIDI notes to sounds.
type SoundMap map[string]*SoundMapEntry

// GetSoundMap fetches the hosted JSON sound map and returns it as a struct.
func GetSoundMap() *SoundMap {
	// Get JSON from Github
	res, err := http.Get(mapURL)
	Must(err)

	// Read response body
	soundMapJSON, err := ioutil.ReadAll(res.Body)
	Must(err)

	// Unmarshal response
	soundMap := &SoundMap{}
	json.Unmarshal(soundMapJSON, soundMap)

	return soundMap
}

// GetEntry retrieves a SoundMap entry given a MIDI note key.
func (sm *SoundMap) GetEntry(key int) (*SoundMapEntry, bool) {
	keyString := strconv.Itoa(key)
	soundMap := *sm
	entry, ok := soundMap[keyString]
	return entry, ok
}

// String prints all contents of the map.
func (sm *SoundMap) String() string {
	s := "[HELP] Currently mapped sounds:\n"
	soundMap := *sm
	for key, sound := range soundMap {
		s += fmt.Sprintf("[HELP]\t\tMIDI Note: %s, Title: %s, Hold To Play: %t, Allow Pausing: %t, Loop: %t\n",
			key, sound.Title, sound.HoldToPlay, sound.AllowPausing, sound.Loop)
	}
	return s
}
