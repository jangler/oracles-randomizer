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

	for sndType, music := range gameMusicData {
		musicMap := make([]mutableRange, 0)
		for i := 0; i < len(music); i++ {
			musicOffset := address{0x39, uint16(soundPointerOffset+3*music[i])}
			musicFullOffset := musicOffset.fullOffset()
			oldData := rom.data[musicFullOffset:musicFullOffset+3]
			musicMap = append(musicMap, mutableRange{musicOffset, oldData, oldData})
		}
		src.Shuffle(len(music), func(i, j int) {
			musicMap[i].new, musicMap[j].new = musicMap[j].new, musicMap[i].new
		})
		for i := 0; i < len(music); i++ {
			rom.codeMutables[fmt.Sprintf("%s%d", sndType, i)] = &musicMap[i]
		}
	}
}