package rom

// treasure sprites
// i could use iota or make a []string instead, but i want explicit numbers.
// the gaps are all green sprites of the left half of impa.
// the rest of the values after the end are just more phonographs.
const (
	spriteRupeeSmall  = 0x00
	spriteRupeeMedium = 0x01
	spriteRupeeLarge  = 0x02
	spriteOreChunk    = 0x03
	spriteLines       = 0x04
	spriteBomb        = 0x05
	spriteEmberSeed   = 0x06
	spriteScentSeed   = 0x07
	spritePegasusSeed = 0x08
	spriteGaleSeed    = 0x09
	spriteMysterySeed = 0x0a
	spriteLines2      = 0x0b
	spriteLines3      = 0x0c
	spriteGashaSeed   = 0x0d
	spriteRing        = 0x0e

	spriteSwordL1         = 0x10
	spriteSwordL2         = 0x11
	spriteSwordL3         = 0x12
	spriteShieldL1        = 0x13
	spriteShieldL2        = 0x14
	spriteShieldL3        = 0x15
	spriteFeatherL1       = 0x16
	spriteFeatherL2       = 0x17
	spriteMagnetGlove     = 0x18
	spriteBracelet        = 0x19
	spriteOreChunk2       = 0x1a
	spriteShovel          = 0x1b
	spriteBoomerangL1     = 0x1c
	spriteBoomerangL2     = 0x1d
	spriteRod             = 0x1e
	spriteSwitchHook      = 0x1f
	spriteSatchel         = 0x20
	spriteSlingshotL1     = 0x21
	spriteSlingshotL2     = 0x22
	spriteStrangeFlute    = 0x23
	spriteBombchu         = 0x24
	spriteBiggoronSword   = 0x25
	spriteMastersPlaque   = 0x26
	spriteRupeeSmallGreen = 0x28

	spriteRupeeSmallCyan    = 0x29
	spriteRupeeSmallYellow  = 0x2a
	spriteRupeeMediumBlue   = 0x2b
	spriteRupeeMediumYellow = 0x2c
	spriteRupeeLargeBlue    = 0x2d
	spriteRupeeLargeYellow  = 0x2e
	spriteOreChunkYellow    = 0x2f
	spriteHalfPotion        = 0x30
	spriteFlippers          = 0x31
	spriteInverseFlipper    = 0x32
	spriteRingBoxL1         = 0x33
	spriteRingBoxL2         = 0x34
	spriteRingBoxL3         = 0x35
	spriteRoundJewel        = 0x36
	spritePyramidJewel      = 0x37
	spriteSquareJewel       = 0x38
	spriteXShapedJewel      = 0x39
	spritePieceOfHeart      = 0x3a
	spriteHeartContainer    = 0x3b

	spriteMap          = 0x40
	spriteCompass      = 0x41
	spriteSmallKey     = 0x42
	spriteBossKey      = 0x43
	spriteGnarledKey   = 0x44
	spriteFloodgateKey = 0x45
	spriteDragonKey    = 0x46
	spriteMakuSeed     = 0x47

	spriteSpringBanana = 0x54
	spriteRickysGloves = 0x55
	spriteBombFlower   = 0x56
	spriteStarOre      = 0x57
	spriteBlueOre      = 0x58
	spriteRedOre       = 0x59
	spriteHardOre      = 0x5a
	spriteRustyBell    = 0x5b
	spritePiratesBell  = 0x5c
	spriteBlueGloves   = 0x5d

	spriteMakuSeed2 = 0x5f
	spriteEssence1  = 0x60
	spriteEssence2  = 0x61
	spriteEssence3  = 0x62
	spriteEssence4  = 0x63
	spriteEssence5  = 0x64
	spriteEssence6  = 0x65
	spriteEssence7  = 0x66
	spriteEssence8  = 0x67

	spriteCuccodex     = 0x70
	spriteLonLonEgg    = 0x71
	spriteGhastlyDoll  = 0x72
	spriteIronPot      = 0x73
	spriteLavaSoup     = 0x74
	spriteGoronVase    = 0x75
	spriteFish         = 0x76
	spriteMegaphone    = 0x77
	spriteMushroom     = 0x78
	spriteWoodenBird   = 0x79
	spriteEngineGrease = 0x7a
	spritePhonograph   = 0x7b
)
