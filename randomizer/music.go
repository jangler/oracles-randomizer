package randomizer

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"math/rand"
)

func (rom *romState) shuffleMusic(src *rand.Rand, game int) {
	soundPointerOffset := sora(game, 0x57cf, 0x5748).(int)

	musicData := make(map[string]map[string][]int)
	if err := yaml.Unmarshal(
		FSMustByte(false, "/romdata/sounds.yaml"), musicData); err != nil {
		panic(err)
	}
	gameMusicData := musicData[gameNames[game]]

	music := gameMusicData["music"]
	musicMap := make([]mutableRange, 0)
	for i := 0; i < len(music); i++ {
		musicOffset := address{0x39, uint16(soundPointerOffset)+3*uint16(music[i])}
		musicFullOffset := musicOffset.fullOffset()
		oldData := rom.data[musicFullOffset:musicFullOffset+3]
		musicMap = append(musicMap, mutableRange{musicOffset, oldData, oldData})
	}
	src.Shuffle(len(music), func(i, j int) {
		tmpRange := musicMap[i].new
		musicMap[i].new = musicMap[j].new
		musicMap[j].new = tmpRange
	})
	for i := 0; i < len(music); i++ {
		rom.codeMutables[fmt.Sprintf("music%d", i)] = &musicMap[i]
	}

	sounds := gameMusicData["sounds"]
	soundsMap := make([]mutableRange, 0)
	for i := 0; i < len(sounds); i++ {
		soundsOffset := address{0x39, uint16(soundPointerOffset)+3*uint16(sounds[i])}
		soundsFullOffset := soundsOffset.fullOffset()
		oldData := rom.data[soundsFullOffset:soundsFullOffset+3]
		soundsMap = append(soundsMap, mutableRange{soundsOffset, oldData, oldData})
	}
	src.Shuffle(len(sounds), func(i, j int) {
		tmpRange := soundsMap[i].new
		soundsMap[i].new = soundsMap[j].new
		soundsMap[j].new = tmpRange
	})
	for i := 0; i < len(sounds); i++ {
		rom.codeMutables[fmt.Sprintf("sounds%d", i)] = &soundsMap[i]
	}
}
