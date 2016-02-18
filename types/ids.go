package types

import "fmt"

// ID represents ID for Minecraft blocks/items.
type ID uint16

// String converts ID to string.
func (id ID) String() string {
	switch id {
	case Air:
		return "Air"
	case Stone:
		return "Stone"
	case Grass:
		return "Grass"
	case Dirt:
		return "Dirt"
	case Cobblestone:
		return "Cobblestone"
	case Plank:
		return "Plank"
	case Sapling:
		return "Sapling"
	case Bedrock:
		return "Bedrock"
	case Water:
		return "Water"
	case StillWater:
		return "StillWater"
	case Lava:
		return "Lava"
	case StillLava:
		return "StillLava"
	case Sand:
		return "Sand"
	case Gravel:
		return "Gravel"
	case GoldOre:
		return "GoldOre"
	case IronOre:
		return "IronOre"
	case CoalOre:
		return "CoalOre"
	case Log:
		return "Log"
	case Leaves:
		return "Leaves"
	case Sponge:
		return "Sponge"
	case Glass:
		return "Glass"
	case LapisOre:
		return "LapisOre"
	case LapisBlock:
		return "LapisBlock"
	case Sandstone:
		return "Sandstone"
	case BedBlock:
		return "BedBlock"
	case Cobweb:
		return "Cobweb"
	case TallGrass:
		return "TallGrass"
	case Bush:
		return "Bush"
	case Wool:
		return "Wool"
	case Dandelion:
		return "Dandelion"
	case Poppy:
		return "Poppy"
	case BrownMushroom:
		return "BrownMushroom"
	case RedMushroom:
		return "RedMushroom"
	case GoldBlock:
		return "GoldBlock"
	case IronBlock:
		return "IronBlock"
	case DoubleSlab:
		return "DoubleSlab"
	case Slab:
		return "Slab"
	case Bricks:
		return "Bricks"
	case Tnt:
		return "Tnt"
	case Bookshelf:
		return "Bookshelf"
	case MossStone:
		return "MossStone"
	case Obsidian:
		return "Obsidian"
	case Torch:
		return "Torch"
	case Fire:
		return "Fire"
	case MonsterSpawner:
		return "MonsterSpawner"
	case WoodStairs:
		return "WoodStairs"
	case Chest:
		return "Chest"
	case DiamondOre:
		return "DiamondOre"
	case DiamondBlock:
		return "DiamondBlock"
	case CraftingTable:
		return "CraftingTable"
	case WheatBlock:
		return "WheatBlock"
	case Farmland:
		return "Farmland"
	case Furnace:
		return "Furnace"
	case BurningFurnace:
		return "BurningFurnace"
	case SignPost:
		return "SignPost"
	case DoorBlock:
		return "DoorBlock"
	case Ladder:
		return "Ladder"
	case CobbleStairs:
		return "CobbleStairs"
	case WallSign:
		return "WallSign"
	case IronDoorBlock:
		return "IronDoorBlock"
	case RedstoneOre:
		return "RedstoneOre"
	case GlowingRedstoneOre:
		return "GlowingRedstoneOre"
	case Snow:
		return "Snow"
	case Ice:
		return "Ice"
	case SnowBlock:
		return "SnowBlock"
	case Cactus:
		return "Cactus"
	case ClayBlock:
		return "ClayBlock"
	case Reeds:
		return "Reeds"
	case Fence:
		return "Fence"
	case Pumpkin:
		return "Pumpkin"
	case Netherrack:
		return "Netherrack"
	case SoulSand:
		return "SoulSand"
	case Glowstone:
		return "Glowstone"
	case LitPumpkin:
		return "LitPumpkin"
	case CakeBlock:
		return "CakeBlock"
	case Trapdoor:
		return "Trapdoor"
	case StoneBricks:
		return "StoneBricks"
	case IronBar:
		return "IronBar"
	case GlassPane:
		return "GlassPane"
	case MelonBlock:
		return "MelonBlock"
	case PumpkinStem:
		return "PumpkinStem"
	case MelonStem:
		return "MelonStem"
	case Vine:
		return "Vine"
	case FenceGate:
		return "FenceGate"
	case BrickStairs:
		return "BrickStairs"
	case StoneBrickStairs:
		return "StoneBrickStairs"
	case Mycelium:
		return "Mycelium"
	case WaterLily:
		return "WaterLily"
	case NetherBricks:
		return "NetherBricks"
	case NetherBrickFence:
		return "NetherBrickFence"
	case NetherBricksStairs:
		return "NetherBricksStairs"
	case EnchantingTable:
		return "EnchantingTable"
	case BrewingStand:
		return "BrewingStand"
	case EndPortal:
		return "EndPortal"
	case EndStone:
		return "EndStone"
	case SandstoneStairs:
		return "SandstoneStairs"
	case EmeraldOre:
		return "EmeraldOre"
	case EmeraldBlock:
		return "EmeraldBlock"
	case SpruceWoodStairs:
		return "SpruceWoodStairs"
	case BirchWoodStairs:
		return "BirchWoodStairs"
	case JungleWoodStairs:
		return "JungleWoodStairs"
	case CobbleWall:
		return "CobbleWall"
	case FlowerPotBlock:
		return "FlowerPotBlock"
	case CarrotBlock:
		return "CarrotBlock"
	case PotatoBlock:
		return "PotatoBlock"
	case Anvil:
		return "Anvil"
	case TrappedChest:
		return "TrappedChest"
	case RedstoneBlock:
		return "RedstoneBlock"
	case QuartzBlock:
		return "QuartzBlock"
	case QuartzStairs:
		return "QuartzStairs"
	case DoubleWoodSlab:
		return "DoubleWoodSlab"
	case WoodSlab:
		return "WoodSlab"
	case StainedClay:
		return "StainedClay"
	case Leaves2:
		return "Leaves2"
	case Wood2:
		return "Wood2"
	case AcaciaWoodStairs:
		return "AcaciaWoodStairs"
	case DarkOakWoodStairs:
		return "DarkOakWoodStairs"
	case IronTrapdoor:
		return "IronTrapdoor"
	case HayBale:
		return "HayBale"
	case Carpet:
		return "Carpet"
	case HardenedClay:
		return "HardenedClay"
	case CoalBlock:
		return "CoalBlock"
	case PackedIce:
		return "PackedIce"
	case DoublePlant:
		return "DoublePlant"
	case FenceGateSpruce:
		return "FenceGateSpruce"
	case FenceGateBirch:
		return "FenceGateBirch"
	case FenceGateJungle:
		return "FenceGateJungle"
	case FenceGateDarkOak:
		return "FenceGateDarkOak"
	case FenceGateAcacia:
		return "FenceGateAcacia"
	case GrassPath:
		return "GrassPath"
	case Podzol:
		return "Podzol"
	case BeetrootBlock:
		return "BeetrootBlock"
	case Stonecutter:
		return "Stonecutter"
	case GlowingObsidian:
		return "GlowingObsidian"
	case IronShovel:
		return "IronShovel"
	case IronPickaxe:
		return "IronPickaxe"
	case IronAxe:
		return "IronAxe"
	case FlintSteel:
		return "FlintSteel"
	case Apple:
		return "Apple"
	case Bow:
		return "Bow"
	case Arrow:
		return "Arrow"
	case Coal:
		return "Coal"
	case Diamond:
		return "Diamond"
	case IronIngot:
		return "IronIngot"
	case GoldIngot:
		return "GoldIngot"
	case IronSword:
		return "IronSword"
	case WoodenSword:
		return "WoodenSword"
	case WoodenShovel:
		return "WoodenShovel"
	case WoodenPickaxe:
		return "WoodenPickaxe"
	case WoodenAxe:
		return "WoodenAxe"
	case StoneSword:
		return "StoneSword"
	case StoneShovel:
		return "StoneShovel"
	case StonePickaxe:
		return "StonePickaxe"
	case StoneAxe:
		return "StoneAxe"
	case DiamondSword:
		return "DiamondSword"
	case DiamondShovel:
		return "DiamondShovel"
	case DiamondPickaxe:
		return "DiamondPickaxe"
	case DiamondAxe:
		return "DiamondAxe"
	case Stick:
		return "Stick"
	case Bowl:
		return "Bowl"
	case MushroomStew:
		return "MushroomStew"
	case GoldSword:
		return "GoldSword"
	case GoldShovel:
		return "GoldShovel"
	case GoldPickaxe:
		return "GoldPickaxe"
	case GoldAxe:
		return "GoldAxe"
	case String:
		return "String"
	case Feather:
		return "Feather"
	case Gunpowder:
		return "Gunpowder"
	case WoodenHoe:
		return "WoodenHoe"
	case StoneHoe:
		return "StoneHoe"
	case IronHoe:
		return "IronHoe"
	case DiamondHoe:
		return "DiamondHoe"
	case GoldHoe:
		return "GoldHoe"
	case Seeds:
		return "Seeds"
	case Wheat:
		return "Wheat"
	case Bread:
		return "Bread"
	case LeatherCap:
		return "LeatherCap"
	case LeatherTunic:
		return "LeatherTunic"
	case LeatherPants:
		return "LeatherPants"
	case LeatherBoots:
		return "LeatherBoots"
	case ChainHelmet:
		return "ChainHelmet"
	case ChainChestplate:
		return "ChainChestplate"
	case ChainLeggings:
		return "ChainLeggings"
	case ChainBoots:
		return "ChainBoots"
	case IronHelmet:
		return "IronHelmet"
	case IronChestplate:
		return "IronChestplate"
	case IronLeggings:
		return "IronLeggings"
	case IronBoots:
		return "IronBoots"
	case DiamondHelmet:
		return "DiamondHelmet"
	case DiamondChestplate:
		return "DiamondChestplate"
	case DiamondLeggings:
		return "DiamondLeggings"
	case DiamondBoots:
		return "DiamondBoots"
	case GoldHelmet:
		return "GoldHelmet"
	case GoldChestplate:
		return "GoldChestplate"
	case GoldLeggings:
		return "GoldLeggings"
	case GoldBoots:
		return "GoldBoots"
	case Flint:
		return "Flint"
	case RawPorkchop:
		return "RawPorkchop"
	case CookedPorkchop:
		return "CookedPorkchop"
	case Painting:
		return "Painting"
	case GoldenApple:
		return "GoldenApple"
	case Sign:
		return "Sign"
	case WoodenDoor:
		return "WoodenDoor"
	case Bucket:
		return "Bucket"
	case Minecart:
		return "Minecart"
	case IronDoor:
		return "IronDoor"
	case Redstone:
		return "Redstone"
	case Snowball:
		return "Snowball"
	case Leather:
		return "Leather"
	case Brick:
		return "Brick"
	case Clay:
		return "Clay"
	case Sugarcane:
		return "Sugarcane"
	case Paper:
		return "Paper"
	case Book:
		return "Book"
	case Slimeball:
		return "Slimeball"
	case Egg:
		return "Egg"
	case Compass:
		return "Compass"
	case FishingRod:
		return "FishingRod"
	case Clock:
		return "Clock"
	case GlowstoneDust:
		return "GlowstoneDust"
	case RawFish:
		return "RawFish"
	case CookedFish:
		return "CookedFish"
	case Dye:
		return "Dye"
	case Bone:
		return "Bone"
	case Sugar:
		return "Sugar"
	case Cake:
		return "Cake"
	case Bed:
		return "Bed"
	case Cookie:
		return "Cookie"
	case Shears:
		return "Shears"
	case Melon:
		return "Melon"
	case PumpkinSeeds:
		return "PumpkinSeeds"
	case MelonSeeds:
		return "MelonSeeds"
	case RawBeef:
		return "RawBeef"
	case Steak:
		return "Steak"
	case RawChicken:
		return "RawChicken"
	case CookedChicken:
		return "CookedChicken"
	case GoldNugget:
		return "GoldNugget"
	case SpawnEgg:
		return "SpawnEgg"
	case Emerald:
		return "Emerald"
	case FlowerPot:
		return "FlowerPot"
	case Carrot:
		return "Carrot"
	case Potato:
		return "Potato"
	case BakedPotato:
		return "BakedPotato"
	case PumpkinPie:
		return "PumpkinPie"
	case NetherBrick:
		return "NetherBrick"
	case Quartz:
		return "Quartz"
	case Camera:
		return "Camera"
	case Beetroot:
		return "Beetroot"
	case BeetrootSeeds:
		return "BeetrootSeeds"
	case BeetrootSoup:
		return "BeetrootSoup"
	default:
		return "Unknown"
	}
}

// Block tries to convert item ID to block ID. If fails, it panics.
func (id ID) Block() byte {
	if id >= 256 {
		panic(fmt.Sprintf("item ID %d(%s) overflows byte", uint16(id), id))
	}
	return byte(id)
}

// item/block IDs
const (
	Air                ID = iota
	Stone                 // 1
	Grass                 // 2
	Dirt                  // 3
	Cobblestone           // 4
	Plank                 // 5
	Sapling               // 6
	Bedrock               // 7
	Water                 // 8
	StillWater            // 9
	Lava                  // 10
	StillLava             // 11
	Sand                  // 12
	Gravel                // 13
	GoldOre               // 14
	IronOre               // 15
	CoalOre               // 16
	Log                   // 17
	Leaves                // 18
	Sponge                // 19
	Glass                 // 20
	LapisOre              // 21
	LapisBlock            // 22
	_                     // 23
	Sandstone             // 24
	_                     // 25
	BedBlock              // 26
	_                     // 27
	_                     // 28
	_                     // 29
	Cobweb                // 30
	TallGrass             // 31
	Bush                  // 32
	_                     // 33
	_                     // 34
	Wool                  // 35
	_                     // 36
	Dandelion             // 37
	Poppy                 // 38
	BrownMushroom         // 39
	RedMushroom           // 40
	GoldBlock             // 41
	IronBlock             // 42
	DoubleSlab            // 43
	Slab                  // 44
	Bricks                // 45
	Tnt                   // 46
	Bookshelf             // 47
	MossStone             // 48
	Obsidian              // 49
	Torch                 // 50
	Fire                  // 51
	MonsterSpawner        // 52
	WoodStairs            // 53
	Chest                 // 54
	_                     // 55
	DiamondOre            // 56
	DiamondBlock          // 57
	CraftingTable         // 58
	WheatBlock            // 59
	Farmland              // 60
	Furnace               // 61
	BurningFurnace        // 62
	SignPost              // 63
	DoorBlock             // 64
	Ladder                // 65
	_                     // 66
	CobbleStairs          // 67
	WallSign              // 68
	_                     // 69
	_                     // 70
	IronDoorBlock         // 71
	_                     // 72
	RedstoneOre           // 73
	GlowingRedstoneOre    // 74
	_                     // 75
	_                     // 76
	_                     // 77
	Snow                  // 78
	Ice                   // 79
	SnowBlock             // 80
	Cactus                // 81
	ClayBlock             // 82
	Reeds                 // 83
	_                     // 84
	Fence                 // 85
	Pumpkin               // 86
	Netherrack            // 87
	SoulSand              // 88
	Glowstone             // 89
	_                     // 90
	LitPumpkin            // 91
	CakeBlock             // 92
	_                     // 93
	_                     // 94
	_                     // 95
	Trapdoor              // 96
	_                     // 97
	StoneBricks           // 98
	_                     // 99
	_                     // 100
	IronBar               // 101
	GlassPane             // 102
	MelonBlock            // 103
	PumpkinStem           // 104
	MelonStem             // 105
	Vine                  // 106
	FenceGate             // 107
	BrickStairs           // 108
	StoneBrickStairs      // 109
	Mycelium              // 110
	WaterLily             // 111
	NetherBricks          // 112
	NetherBrickFence      // 113
	NetherBricksStairs    // 114
	_                     // 115
	EnchantingTable       // 116
	BrewingStand          // 117
	_                     // 118
	_                     // 119
	EndPortal             // 120
	EndStone              // 121
	_                     // 122
	_                     // 123
	_                     // 124
	_                     // 125
	_                     // 126
	_                     // 127
	SandstoneStairs       // 128
	EmeraldOre            // 129
	_                     // 130
	_                     // 131
	_                     // 132
	EmeraldBlock          // 133
	SpruceWoodStairs      // 134
	BirchWoodStairs       // 135
	JungleWoodStairs      // 136
	_                     // 137
	_                     // 138
	CobbleWall            // 139
	FlowerPotBlock        // 140
	CarrotBlock           // 141
	PotatoBlock           // 142
	_                     // 143
	_                     // 144
	Anvil                 // 145
	TrappedChest          // 146
	_                     // 147
	_                     // 148
	_                     // 149
	_                     // 150
	_                     // 151
	RedstoneBlock         // 152
	_                     // 153
	_                     // 154
	QuartzBlock           // 155
	QuartzStairs          // 156
	DoubleWoodSlab        // 157
	WoodSlab              // 158
	StainedClay           // 159
	_                     // 160
	Leaves2               // 161
	Wood2                 // 162
	AcaciaWoodStairs      // 163
	DarkOakWoodStairs     // 164
	_                     // 165
	_                     // 166
	IronTrapdoor          // 167
	_                     // 168
	_                     // 169
	HayBale               // 170
	Carpet                // 171
	HardenedClay          // 172
	CoalBlock             // 173
	PackedIce             // 174
	DoublePlant           // 175
	_                     // 176
	_                     // 177
	_                     // 178
	_                     // 179
	_                     // 180
	_                     // 181
	_                     // 182
	FenceGateSpruce       // 183
	FenceGateBirch        // 184
	FenceGateJungle       // 185
	FenceGateDarkOak      // 186
	FenceGateAcacia       // 187
	_                     // 188
	_                     // 189
	_                     // 190
	_                     // 191
	_                     // 192
	_                     // 193
	_                     // 194
	_                     // 195
	_                     // 196
	_                     // 197
	GrassPath             // 198
	_                     // 199
	_                     // 200
	_                     // 201
	_                     // 202
	_                     // 203
	_                     // 204
	_                     // 205
	_                     // 206
	_                     // 207
	_                     // 208
	_                     // 209
	_                     // 210
	_                     // 211
	_                     // 212
	_                     // 213
	_                     // 214
	_                     // 215
	_                     // 216
	_                     // 217
	_                     // 218
	_                     // 219
	_                     // 220
	_                     // 221
	_                     // 222
	_                     // 223
	_                     // 224
	_                     // 225
	_                     // 226
	_                     // 227
	_                     // 228
	_                     // 229
	_                     // 230
	_                     // 231
	_                     // 232
	_                     // 233
	_                     // 234
	_                     // 235
	_                     // 236
	_                     // 237
	_                     // 238
	_                     // 239
	_                     // 240
	_                     // 241
	_                     // 242
	Podzol                // 243
	BeetrootBlock         // 244
	Stonecutter           // 245
	GlowingObsidian       // 246
	_                     // 247
	_                     // 248
	_                     // 249
	_                     // 250
	_                     // 251
	_                     // 252
	_                     // 253
	_                     // 254
	_                     // 255
	IronShovel            // 256
	IronPickaxe           // 257
	IronAxe               // 258
	FlintSteel            // 259
	Apple                 // 260
	Bow                   // 261
	Arrow                 // 262
	Coal                  // 263
	Diamond               // 264
	IronIngot             // 265
	GoldIngot             // 266
	IronSword             // 267
	WoodenSword           // 268
	WoodenShovel          // 269
	WoodenPickaxe         // 270
	WoodenAxe             // 271
	StoneSword            // 272
	StoneShovel           // 273
	StonePickaxe          // 274
	StoneAxe              // 275
	DiamondSword          // 276
	DiamondShovel         // 277
	DiamondPickaxe        // 278
	DiamondAxe            // 279
	Stick                 // 280
	Bowl                  // 281
	MushroomStew          // 282
	GoldSword             // 283
	GoldShovel            // 284
	GoldPickaxe           // 285
	GoldAxe               // 286
	String                // 287
	Feather               // 288
	Gunpowder             // 289
	WoodenHoe             // 290
	StoneHoe              // 291
	IronHoe               // 292
	DiamondHoe            // 293
	GoldHoe               // 294
	Seeds                 // 295
	Wheat                 // 296
	Bread                 // 297
	LeatherCap            // 298
	LeatherTunic          // 299
	LeatherPants          // 300
	LeatherBoots          // 301
	ChainHelmet           // 302
	ChainChestplate       // 303
	ChainLeggings         // 304
	ChainBoots            // 305
	IronHelmet            // 306
	IronChestplate        // 307
	IronLeggings          // 308
	IronBoots             // 309
	DiamondHelmet         // 310
	DiamondChestplate     // 311
	DiamondLeggings       // 312
	DiamondBoots          // 313
	GoldHelmet            // 314
	GoldChestplate        // 315
	GoldLeggings          // 316
	GoldBoots             // 317
	Flint                 // 318
	RawPorkchop           // 319
	CookedPorkchop        // 320
	Painting              // 321
	GoldenApple           // 322
	Sign                  // 323
	WoodenDoor            // 324
	Bucket                // 325
	_                     // 326
	_                     // 327
	Minecart              // 328
	_                     // 329
	IronDoor              // 330
	Redstone              // 331
	Snowball              // 332
	_                     // 333
	Leather               // 334
	_                     // 335
	Brick                 // 336
	Clay                  // 337
	Sugarcane             // 338
	Paper                 // 339
	Book                  // 340
	Slimeball             // 341
	_                     // 342
	_                     // 343
	Egg                   // 344
	Compass               // 345
	FishingRod            // 346
	Clock                 // 347
	GlowstoneDust         // 348
	RawFish               // 349
	CookedFish            // 350
	Dye                   // 351
	Bone                  // 352
	Sugar                 // 353
	Cake                  // 354
	Bed                   // 355
	_                     // 356
	Cookie                // 357
	_                     // 358
	Shears                // 359
	Melon                 // 360
	PumpkinSeeds          // 361
	MelonSeeds            // 362
	RawBeef               // 363
	Steak                 // 364
	RawChicken            // 365
	CookedChicken         // 366
	_                     // 367
	_                     // 368
	_                     // 369
	_                     // 370
	GoldNugget            // 371
	_                     // 372
	_                     // 373
	_                     // 374
	_                     // 375
	_                     // 376
	_                     // 377
	_                     // 378
	_                     // 379
	_                     // 380
	_                     // 381
	_                     // 382
	SpawnEgg              // 383
	_                     // 384
	_                     // 385
	_                     // 386
	_                     // 387
	Emerald               // 388
	_                     // 389
	FlowerPot             // 390
	Carrot                // 391
	Potato                // 392
	BakedPotato           // 393
	_                     // 394
	_                     // 395
	_                     // 396
	_                     // 397
	_                     // 398
	_                     // 399
	PumpkinPie            // 400
	_                     // 401
	_                     // 402
	_                     // 403
	_                     // 404
	NetherBrick           // 405
	Quartz                // 406
	_                     // 407
	_                     // 408
	_                     // 409
	_                     // 410
	_                     // 411
	_                     // 412
	_                     // 413
	_                     // 414
	_                     // 415
	_                     // 416
	_                     // 417
	_                     // 418
	_                     // 419
	_                     // 420
	_                     // 421
	_                     // 422
	_                     // 423
	_                     // 424
	_                     // 425
	_                     // 426
	_                     // 427
	_                     // 428
	_                     // 429
	_                     // 430
	_                     // 431
	_                     // 432
	_                     // 433
	_                     // 434
	_                     // 435
	_                     // 436
	_                     // 437
	_                     // 438
	_                     // 439
	_                     // 440
	_                     // 441
	_                     // 442
	_                     // 443
	_                     // 444
	_                     // 445
	_                     // 446
	_                     // 447
	_                     // 448
	_                     // 449
	_                     // 450
	_                     // 451
	_                     // 452
	_                     // 453
	_                     // 454
	_                     // 455
	Camera                // 456
	Beetroot              // 457
	BeetrootSeeds         // 458
	BeetrootSoup          // 459
)

// aliases
const (
	Carrots             = Carrot             // 391
	Slabs               = Slab               // 44
	NetherBrickBlock    = NetherBricks       // 112
	EnchantTable        = EnchantingTable    // 116
	Planks              = Plank              // 5
	JackOLantern        = LitPumpkin         // 91
	SnowLayer           = Snow               // 78
	DarkOakWoodenStairs = DarkOakWoodStairs  // 164
	GoldenHoe           = GoldHoe            // 294
	LitRedstoneOre      = GlowingRedstoneOre // 74
	Leave               = Leaves             // 18
	CobblestoneStairs   = CobbleStairs       // 67
	OakWoodenStairs     = WoodStairs         // 53
	StoneWall           = CobbleWall         // 139
	JungleWoodenStairs  = JungleWoodStairs   // 136
	GoldenSword         = GoldSword          // 283
	StainedHardenedClay = StainedClay        // 159
	LitFurnace          = BurningFurnace     // 62
	MossyStone          = MossStone          // 48
	WoodSlabs           = WoodSlab           // 158
	IronBars            = IronBar            // 101
	Rose                = Poppy              // 38
	WoodenStairs        = WoodStairs         // 53
	BeetrootSeed        = BeetrootSeeds      // 458
	Potatoes            = Potato             // 392
	EnchantmentTable    = EnchantingTable    // 116
	Leave2              = Leaves2            // 161
	BirchWoodenStairs   = BirchWoodStairs    // 135
	Log2                = Wood2              // 162
	StoneBrick          = StoneBricks        // 98
	DoubleWoodenSlabs   = DoubleWoodSlab     // 157
	Wood                = Log                // 17
	Saplings            = Sapling            // 6
	GlassPanel          = GlassPane          // 102
	FlintAndSteel       = FlintSteel         // 259
	DoubleWoodenSlab    = DoubleWoodSlab     // 157
	MelonSlice          = Melon              // 360
	Sticks              = Stick              // 280
	RedstoneDust        = Redstone           // 331
	WoodenPlank         = Plank              // 5
	GlowstoneBlock      = Glowstone          // 89
	AcaciaWoodenStairs  = AcaciaWoodStairs   // 163
	Trunk2              = Wood2              // 162
	WheatSeeds          = Seeds              // 295
	CobblestoneWall     = CobbleWall         // 139
	DoubleSlabs         = DoubleSlab         // 43
	SpruceWoodenStairs  = SpruceWoodStairs   // 134
	Vines               = Vine               // 106
	GoldenNugget        = GoldNugget         // 371
	WoodenSlab          = WoodSlab           // 158
	BricksBlock         = Bricks             // 45
	SugarCane           = Sugarcane          // 338
	Workbench           = CraftingTable      // 58
	LilyPad             = WaterLily          // 111
	SugarCanes          = Sugarcane          // 338
	SugarcaneBlock      = Reeds              // 83
	GoldenAxe           = GoldAxe            // 286
	NetherQuartz        = Quartz             // 406
	BakedPotatoes       = BakedPotato        // 393
	WoodenPlanks        = Plank              // 5
	WoodenSlabs         = WoodSlab           // 158
	CookedBeef          = Steak              // 364
	DoubleWoodSlabs     = DoubleWoodSlab     // 157
	GoldenShovel        = GoldShovel         // 284
	Trunk               = Log                // 17
	RedFlower           = Poppy              // 38
	WoodenDoorBlock     = DoorBlock          // 64
	WoodDoorBlock       = DoorBlock          // 64
	GoldenPickaxe       = GoldPickaxe        // 285
	OakWoodStairs       = WoodStairs         // 53
	Cobble              = Cobblestone        // 4
	DeadBush            = Bush               // 32
)

var idMap = map[string]ID{
	"Stone":              Stone,              // 1
	"Grass":              Grass,              // 2
	"Dirt":               Dirt,               // 3
	"Cobblestone":        Cobblestone,        // 4
	"Plank":              Plank,              // 5
	"Sapling":            Sapling,            // 6
	"Bedrock":            Bedrock,            // 7
	"Water":              Water,              // 8
	"StillWater":         StillWater,         // 9
	"Lava":               Lava,               // 10
	"StillLava":          StillLava,          // 11
	"Sand":               Sand,               // 12
	"Gravel":             Gravel,             // 13
	"GoldOre":            GoldOre,            // 14
	"IronOre":            IronOre,            // 15
	"CoalOre":            CoalOre,            // 16
	"Log":                Log,                // 17
	"Leaves":             Leaves,             // 18
	"Sponge":             Sponge,             // 19
	"Glass":              Glass,              // 20
	"LapisOre":           LapisOre,           // 21
	"LapisBlock":         LapisBlock,         // 22
	"Sandstone":          Sandstone,          // 24
	"BedBlock":           BedBlock,           // 26
	"Cobweb":             Cobweb,             // 30
	"TallGrass":          TallGrass,          // 31
	"Bush":               Bush,               // 32
	"Wool":               Wool,               // 35
	"Dandelion":          Dandelion,          // 37
	"Poppy":              Poppy,              // 38
	"BrownMushroom":      BrownMushroom,      // 39
	"RedMushroom":        RedMushroom,        // 40
	"GoldBlock":          GoldBlock,          // 41
	"IronBlock":          IronBlock,          // 42
	"DoubleSlab":         DoubleSlab,         // 43
	"Slab":               Slab,               // 44
	"Bricks":             Bricks,             // 45
	"Tnt":                Tnt,                // 46
	"Bookshelf":          Bookshelf,          // 47
	"MossStone":          MossStone,          // 48
	"Obsidian":           Obsidian,           // 49
	"Torch":              Torch,              // 50
	"Fire":               Fire,               // 51
	"MonsterSpawner":     MonsterSpawner,     // 52
	"WoodStairs":         WoodStairs,         // 53
	"Chest":              Chest,              // 54
	"DiamondOre":         DiamondOre,         // 56
	"DiamondBlock":       DiamondBlock,       // 57
	"CraftingTable":      CraftingTable,      // 58
	"WheatBlock":         WheatBlock,         // 59
	"Farmland":           Farmland,           // 60
	"Furnace":            Furnace,            // 61
	"BurningFurnace":     BurningFurnace,     // 62
	"SignPost":           SignPost,           // 63
	"DoorBlock":          DoorBlock,          // 64
	"Ladder":             Ladder,             // 65
	"CobbleStairs":       CobbleStairs,       // 67
	"WallSign":           WallSign,           // 68
	"IronDoorBlock":      IronDoorBlock,      // 71
	"RedstoneOre":        RedstoneOre,        // 73
	"GlowingRedstoneOre": GlowingRedstoneOre, // 74
	"Snow":               Snow,               // 78
	"Ice":                Ice,                // 79
	"SnowBlock":          SnowBlock,          // 80
	"Cactus":             Cactus,             // 81
	"ClayBlock":          ClayBlock,          // 82
	"Reeds":              Reeds,              // 83
	"Fence":              Fence,              // 85
	"Pumpkin":            Pumpkin,            // 86
	"Netherrack":         Netherrack,         // 87
	"SoulSand":           SoulSand,           // 88
	"Glowstone":          Glowstone,          // 89
	"LitPumpkin":         LitPumpkin,         // 91
	"CakeBlock":          CakeBlock,          // 92
	"Trapdoor":           Trapdoor,           // 96
	"StoneBricks":        StoneBricks,        // 98
	"IronBar":            IronBar,            // 101
	"GlassPane":          GlassPane,          // 102
	"MelonBlock":         MelonBlock,         // 103
	"PumpkinStem":        PumpkinStem,        // 104
	"MelonStem":          MelonStem,          // 105
	"Vine":               Vine,               // 106
	"FenceGate":          FenceGate,          // 107
	"BrickStairs":        BrickStairs,        // 108
	"StoneBrickStairs":   StoneBrickStairs,   // 109
	"Mycelium":           Mycelium,           // 110
	"WaterLily":          WaterLily,          // 111
	"NetherBricks":       NetherBricks,       // 112
	"NetherBrickFence":   NetherBrickFence,   // 113
	"NetherBricksStairs": NetherBricksStairs, // 114
	"EnchantingTable":    EnchantingTable,    // 116
	"BrewingStand":       BrewingStand,       // 117
	"EndPortal":          EndPortal,          // 120
	"EndStone":           EndStone,           // 121
	"SandstoneStairs":    SandstoneStairs,    // 128
	"EmeraldOre":         EmeraldOre,         // 129
	"EmeraldBlock":       EmeraldBlock,       // 133
	"SpruceWoodStairs":   SpruceWoodStairs,   // 134
	"BirchWoodStairs":    BirchWoodStairs,    // 135
	"JungleWoodStairs":   JungleWoodStairs,   // 136
	"CobbleWall":         CobbleWall,         // 139
	"FlowerPotBlock":     FlowerPotBlock,     // 140
	"CarrotBlock":        CarrotBlock,        // 141
	"PotatoBlock":        PotatoBlock,        // 142
	"Anvil":              Anvil,              // 145
	"TrappedChest":       TrappedChest,       // 146
	"RedstoneBlock":      RedstoneBlock,      // 152
	"QuartzBlock":        QuartzBlock,        // 155
	"QuartzStairs":       QuartzStairs,       // 156
	"DoubleWoodSlab":     DoubleWoodSlab,     // 157
	"WoodSlab":           WoodSlab,           // 158
	"StainedClay":        StainedClay,        // 159
	"Leaves2":            Leaves2,            // 161
	"Wood2":              Wood2,              // 162
	"AcaciaWoodStairs":   AcaciaWoodStairs,   // 163
	"DarkOakWoodStairs":  DarkOakWoodStairs,  // 164
	"IronTrapdoor":       IronTrapdoor,       // 167
	"HayBale":            HayBale,            // 170
	"Carpet":             Carpet,             // 171
	"HardenedClay":       HardenedClay,       // 172
	"CoalBlock":          CoalBlock,          // 173
	"PackedIce":          PackedIce,          // 174
	"DoublePlant":        DoublePlant,        // 175
	"FenceGateSpruce":    FenceGateSpruce,    // 183
	"FenceGateBirch":     FenceGateBirch,     // 184
	"FenceGateJungle":    FenceGateJungle,    // 185
	"FenceGateDarkOak":   FenceGateDarkOak,   // 186
	"FenceGateAcacia":    FenceGateAcacia,    // 187
	"GrassPath":          GrassPath,          // 198
	"Podzol":             Podzol,             // 243
	"BeetrootBlock":      BeetrootBlock,      // 244
	"Stonecutter":        Stonecutter,        // 245
	"GlowingObsidian":    GlowingObsidian,    // 246
	"IronShovel":         IronShovel,         // 256
	"IronPickaxe":        IronPickaxe,        // 257
	"IronAxe":            IronAxe,            // 258
	"FlintSteel":         FlintSteel,         // 259
	"Apple":              Apple,              // 260
	"Bow":                Bow,                // 261
	"Arrow":              Arrow,              // 262
	"Coal":               Coal,               // 263
	"Diamond":            Diamond,            // 264
	"IronIngot":          IronIngot,          // 265
	"GoldIngot":          GoldIngot,          // 266
	"IronSword":          IronSword,          // 267
	"WoodenSword":        WoodenSword,        // 268
	"WoodenShovel":       WoodenShovel,       // 269
	"WoodenPickaxe":      WoodenPickaxe,      // 270
	"WoodenAxe":          WoodenAxe,          // 271
	"StoneSword":         StoneSword,         // 272
	"StoneShovel":        StoneShovel,        // 273
	"StonePickaxe":       StonePickaxe,       // 274
	"StoneAxe":           StoneAxe,           // 275
	"DiamondSword":       DiamondSword,       // 276
	"DiamondShovel":      DiamondShovel,      // 277
	"DiamondPickaxe":     DiamondPickaxe,     // 278
	"DiamondAxe":         DiamondAxe,         // 279
	"Stick":              Stick,              // 280
	"Bowl":               Bowl,               // 281
	"MushroomStew":       MushroomStew,       // 282
	"GoldSword":          GoldSword,          // 283
	"GoldShovel":         GoldShovel,         // 284
	"GoldPickaxe":        GoldPickaxe,        // 285
	"GoldAxe":            GoldAxe,            // 286
	"String":             String,             // 287
	"Feather":            Feather,            // 288
	"Gunpowder":          Gunpowder,          // 289
	"WoodenHoe":          WoodenHoe,          // 290
	"StoneHoe":           StoneHoe,           // 291
	"IronHoe":            IronHoe,            // 292
	"DiamondHoe":         DiamondHoe,         // 293
	"GoldHoe":            GoldHoe,            // 294
	"Seeds":              Seeds,              // 295
	"Wheat":              Wheat,              // 296
	"Bread":              Bread,              // 297
	"LeatherCap":         LeatherCap,         // 298
	"LeatherTunic":       LeatherTunic,       // 299
	"LeatherPants":       LeatherPants,       // 300
	"LeatherBoots":       LeatherBoots,       // 301
	"ChainHelmet":        ChainHelmet,        // 302
	"ChainChestplate":    ChainChestplate,    // 303
	"ChainLeggings":      ChainLeggings,      // 304
	"ChainBoots":         ChainBoots,         // 305
	"IronHelmet":         IronHelmet,         // 306
	"IronChestplate":     IronChestplate,     // 307
	"IronLeggings":       IronLeggings,       // 308
	"IronBoots":          IronBoots,          // 309
	"DiamondHelmet":      DiamondHelmet,      // 310
	"DiamondChestplate":  DiamondChestplate,  // 311
	"DiamondLeggings":    DiamondLeggings,    // 312
	"DiamondBoots":       DiamondBoots,       // 313
	"GoldHelmet":         GoldHelmet,         // 314
	"GoldChestplate":     GoldChestplate,     // 315
	"GoldLeggings":       GoldLeggings,       // 316
	"GoldBoots":          GoldBoots,          // 317
	"Flint":              Flint,              // 318
	"RawPorkchop":        RawPorkchop,        // 319
	"CookedPorkchop":     CookedPorkchop,     // 320
	"Painting":           Painting,           // 321
	"GoldenApple":        GoldenApple,        // 322
	"Sign":               Sign,               // 323
	"WoodenDoor":         WoodenDoor,         // 324
	"Bucket":             Bucket,             // 325
	"Minecart":           Minecart,           // 328
	"IronDoor":           IronDoor,           // 330
	"Redstone":           Redstone,           // 331
	"Snowball":           Snowball,           // 332
	"Leather":            Leather,            // 334
	"Brick":              Brick,              // 336
	"Clay":               Clay,               // 337
	"Sugarcane":          Sugarcane,          // 338
	"Paper":              Paper,              // 339
	"Book":               Book,               // 340
	"Slimeball":          Slimeball,          // 341
	"Egg":                Egg,                // 344
	"Compass":            Compass,            // 345
	"FishingRod":         FishingRod,         // 346
	"Clock":              Clock,              // 347
	"GlowstoneDust":      GlowstoneDust,      // 348
	"RawFish":            RawFish,            // 349
	"CookedFish":         CookedFish,         // 350
	"Dye":                Dye,                // 351
	"Bone":               Bone,               // 352
	"Sugar":              Sugar,              // 353
	"Cake":               Cake,               // 354
	"Bed":                Bed,                // 355
	"Cookie":             Cookie,             // 357
	"Shears":             Shears,             // 359
	"Melon":              Melon,              // 360
	"PumpkinSeeds":       PumpkinSeeds,       // 361
	"MelonSeeds":         MelonSeeds,         // 362
	"RawBeef":            RawBeef,            // 363
	"Steak":              Steak,              // 364
	"RawChicken":         RawChicken,         // 365
	"CookedChicken":      CookedChicken,      // 366
	"GoldNugget":         GoldNugget,         // 371
	"SpawnEgg":           SpawnEgg,           // 383
	"Emerald":            Emerald,            // 388
	"FlowerPot":          FlowerPot,          // 390
	"Carrot":             Carrot,             // 391
	"Potato":             Potato,             // 392
	"BakedPotato":        BakedPotato,        // 393
	"PumpkinPie":         PumpkinPie,         // 400
	"NetherBrick":        NetherBrick,        // 405
	"Quartz":             Quartz,             // 406
	"Camera":             Camera,             // 456
	"Beetroot":           Beetroot,           // 457
	"BeetrootSeeds":      BeetrootSeeds,      // 458
	"BeetrootSoup":       BeetrootSoup,       // 459

	//aliases

	"Carrots":             Carrots,             // 391
	"Slabs":               Slabs,               // 44
	"NetherBrickBlock":    NetherBrickBlock,    // 112
	"EnchantTable":        EnchantTable,        // 116
	"Planks":              Planks,              // 5
	"JackOLantern":        JackOLantern,        // 91
	"SnowLayer":           SnowLayer,           // 78
	"DarkOakWoodenStairs": DarkOakWoodenStairs, // 164
	"GoldenHoe":           GoldenHoe,           // 294
	"LitRedstoneOre":      LitRedstoneOre,      // 74
	"Leave":               Leave,               // 18
	"CobblestoneStairs":   CobblestoneStairs,   // 67
	"OakWoodenStairs":     OakWoodenStairs,     // 53
	"StoneWall":           StoneWall,           // 139
	"JungleWoodenStairs":  JungleWoodenStairs,  // 136
	"GoldenSword":         GoldenSword,         // 283
	"StainedHardenedClay": StainedHardenedClay, // 159
	"LitFurnace":          LitFurnace,          // 62
	"MossyStone":          MossyStone,          // 48
	"WoodSlabs":           WoodSlabs,           // 158
	"IronBars":            IronBars,            // 101
	"Rose":                Rose,                // 38
	"WoodenStairs":        WoodenStairs,        // 53
	"BeetrootSeed":        BeetrootSeed,        // 458
	"Potatoes":            Potatoes,            // 392
	"EnchantmentTable":    EnchantmentTable,    // 116
	"Leave2":              Leave2,              // 161
	"BirchWoodenStairs":   BirchWoodenStairs,   // 135
	"Log2":                Log2,                // 162
	"StoneBrick":          StoneBrick,          // 98
	"DoubleWoodenSlabs":   DoubleWoodenSlabs,   // 157
	"Wood":                Wood,                // 17
	"Saplings":            Saplings,            // 6
	"GlassPanel":          GlassPanel,          // 102
	"FlintAndSteel":       FlintAndSteel,       // 259
	"DoubleWoodenSlab":    DoubleWoodenSlab,    // 157
	"MelonSlice":          MelonSlice,          // 360
	"Sticks":              Sticks,              // 280
	"RedstoneDust":        RedstoneDust,        // 331
	"WoodenPlank":         WoodenPlank,         // 5
	"GlowstoneBlock":      GlowstoneBlock,      // 89
	"AcaciaWoodenStairs":  AcaciaWoodenStairs,  // 163
	"Trunk2":              Trunk2,              // 162
	"WheatSeeds":          WheatSeeds,          // 295
	"CobblestoneWall":     CobblestoneWall,     // 139
	"DoubleSlabs":         DoubleSlabs,         // 43
	"SpruceWoodenStairs":  SpruceWoodenStairs,  // 134
	"Vines":               Vines,               // 106
	"GoldenNugget":        GoldenNugget,        // 371
	"WoodenSlab":          WoodenSlab,          // 158
	"BricksBlock":         BricksBlock,         // 45
	"SugarCane":           SugarCane,           // 338
	"Workbench":           Workbench,           // 58
	"LilyPad":             LilyPad,             // 111
	"SugarCanes":          SugarCanes,          // 338
	"SugarcaneBlock":      SugarcaneBlock,      // 83
	"GoldenAxe":           GoldenAxe,           // 286
	"NetherQuartz":        NetherQuartz,        // 406
	"BakedPotatoes":       BakedPotatoes,       // 393
	"WoodenPlanks":        WoodenPlanks,        // 5
	"WoodenSlabs":         WoodenSlabs,         // 158
	"CookedBeef":          CookedBeef,          // 364
	"DoubleWoodSlabs":     DoubleWoodSlabs,     // 157
	"GoldenShovel":        GoldenShovel,        // 284
	"Trunk":               Trunk,               // 17
	"RedFlower":           RedFlower,           // 38
	"WoodenDoorBlock":     WoodenDoorBlock,     // 64
	"WoodDoorBlock":       WoodDoorBlock,       // 64
	"GoldenPickaxe":       GoldenPickaxe,       // 285
	"OakWoodStairs":       OakWoodStairs,       // 53
	"Cobble":              Cobble,              // 4
	"DeadBush":            DeadBush,            // 32
}

func StringID(name string) ID {
	if id, ok := idMap[name]; ok {
		return id
	}
	return 65535
}
