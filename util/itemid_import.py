# Imports PocketMine-MP's item ID constant declaration code
import re

toimport = '''
	//All Block IDs are here too
	const AIR = 0;
	const STONE = 1;
	const GRASS = 2;
	const DIRT = 3;
	const COBBLESTONE = 4;
	const COBBLE = 4;
	const PLANK = 5;
	const PLANKS = 5;
	const WOODEN_PLANK = 5;
	const WOODEN_PLANKS = 5;
	const SAPLING = 6;
	const SAPLINGS = 6;
	const BEDROCK = 7;
	const WATER = 8;
	const STILL_WATER = 9;
	const LAVA = 10;
	const STILL_LAVA = 11;
	const SAND = 12;
	const GRAVEL = 13;
	const GOLD_ORE = 14;
	const IRON_ORE = 15;
	const COAL_ORE = 16;
	const LOG = 17;
	const WOOD = 17;
	const TRUNK = 17;
	const LEAVES = 18;
	const LEAVE = 18;
	const SPONGE = 19;
	const GLASS = 20;
	const LAPIS_ORE = 21;
	const LAPIS_BLOCK = 22;
	const SANDSTONE = 24;
	const BED_BLOCK = 26;
	const COBWEB = 30;
	const TALL_GRASS = 31;
	const BUSH = 32;
	const DEAD_BUSH = 32;
	const WOOL = 35;
	const DANDELION = 37;
	const POPPY = 38;
	const ROSE = 38;
	const RED_FLOWER = 38;
	const BROWN_MUSHROOM = 39;
	const RED_MUSHROOM = 40;
	const GOLD_BLOCK = 41;
	const IRON_BLOCK = 42;
	const DOUBLE_SLAB = 43;
	const DOUBLE_SLABS = 43;
	const SLAB = 44;
	const SLABS = 44;
	const BRICKS = 45;
	const BRICKS_BLOCK = 45;
	const TNT = 46;
	const BOOKSHELF = 47;
	const MOSS_STONE = 48;
	const MOSSY_STONE = 48;
	const OBSIDIAN = 49;
	const TORCH = 50;
	const FIRE = 51;
	const MONSTER_SPAWNER = 52;
	const WOOD_STAIRS = 53;
	const WOODEN_STAIRS = 53;
	const OAK_WOOD_STAIRS = 53;
	const OAK_WOODEN_STAIRS = 53;
	const CHEST = 54;
	const DIAMOND_ORE = 56;
	const DIAMOND_BLOCK = 57;
	const CRAFTING_TABLE = 58;
	const WORKBENCH = 58;
	const WHEAT_BLOCK = 59;
	const FARMLAND = 60;
	const FURNACE = 61;
	const BURNING_FURNACE = 62;
	const LIT_FURNACE = 62;
	const SIGN_POST = 63;
	const DOOR_BLOCK = 64;
	const WOODEN_DOOR_BLOCK = 64;
	const WOOD_DOOR_BLOCK = 64;
	const LADDER = 65;
	const COBBLE_STAIRS = 67;
	const COBBLESTONE_STAIRS = 67;
	const WALL_SIGN = 68;
	const IRON_DOOR_BLOCK = 71;
	const REDSTONE_ORE = 73;
	const GLOWING_REDSTONE_ORE = 74;
	const LIT_REDSTONE_ORE = 74;
	const SNOW = 78;
	const SNOW_LAYER = 78;
	const ICE = 79;
	const SNOW_BLOCK = 80;
	const CACTUS = 81;
	const CLAY_BLOCK = 82;
	const REEDS = 83;
	const SUGARCANE_BLOCK = 83;
	const FENCE = 85;
	const PUMPKIN = 86;
	const NETHERRACK = 87;
	const SOUL_SAND = 88;
	const GLOWSTONE = 89;
	const GLOWSTONE_BLOCK = 89;
	const LIT_PUMPKIN = 91;
	const JACK_O_LANTERN = 91;
	const CAKE_BLOCK = 92;
	const TRAPDOOR = 96;
	const STONE_BRICKS = 98;
	const STONE_BRICK = 98;
	const IRON_BAR = 101;
	const IRON_BARS = 101;
	const GLASS_PANE = 102;
	const GLASS_PANEL = 102;
	const MELON_BLOCK = 103;
	const PUMPKIN_STEM = 104;
	const MELON_STEM = 105;
	const VINE = 106;
	const VINES = 106;
	const FENCE_GATE = 107;
	const BRICK_STAIRS = 108;
	const STONE_BRICK_STAIRS = 109;
	const MYCELIUM = 110;
	const WATER_LILY = 111;
	const LILY_PAD = 111;
	const NETHER_BRICKS = 112;
	const NETHER_BRICK_BLOCK = 112;
	const NETHER_BRICK_FENCE = 113;
	const NETHER_BRICKS_STAIRS = 114;
	const ENCHANTING_TABLE = 116;
	const ENCHANT_TABLE = 116;
	const ENCHANTMENT_TABLE = 116;
	const BREWING_STAND = 117;
	const END_PORTAL = 120;
	const END_STONE = 121;
	const SANDSTONE_STAIRS = 128;
	const EMERALD_ORE = 129;
	const EMERALD_BLOCK = 133;
	const SPRUCE_WOOD_STAIRS = 134;
	const SPRUCE_WOODEN_STAIRS = 134;
	const BIRCH_WOOD_STAIRS = 135;
	const BIRCH_WOODEN_STAIRS = 135;
	const JUNGLE_WOOD_STAIRS = 136;
	const JUNGLE_WOODEN_STAIRS = 136;
	const COBBLE_WALL = 139;
	const STONE_WALL = 139;
	const COBBLESTONE_WALL = 139;
	const FLOWER_POT_BLOCK = 140;
	const CARROT_BLOCK = 141;
	const POTATO_BLOCK = 142;
	const ANVIL = 145;
	const TRAPPED_CHEST = 146;
	const REDSTONE_BLOCK = 152;
	const QUARTZ_BLOCK = 155;
	const QUARTZ_STAIRS = 156;
	const DOUBLE_WOOD_SLAB = 157;
	const DOUBLE_WOODEN_SLAB = 157;
	const DOUBLE_WOOD_SLABS = 157;
	const DOUBLE_WOODEN_SLABS = 157;
	const WOOD_SLAB = 158;
	const WOODEN_SLAB = 158;
	const WOOD_SLABS = 158;
	const WOODEN_SLABS = 158;
	const STAINED_CLAY = 159;
	const STAINED_HARDENED_CLAY = 159;
	const LEAVES2 = 161;
	const LEAVE2 = 161;
	const WOOD2 = 162;
	const TRUNK2 = 162;
	const LOG2 = 162;
	const ACACIA_WOOD_STAIRS = 163;
	const ACACIA_WOODEN_STAIRS = 163;
	const DARK_OAK_WOOD_STAIRS = 164;
	const DARK_OAK_WOODEN_STAIRS = 164;
	const IRON_TRAPDOOR = 167;
	const HAY_BALE = 170;
	const CARPET = 171;
	const HARDENED_CLAY = 172;
	const COAL_BLOCK = 173;
	const PACKED_ICE = 174;
	const DOUBLE_PLANT = 175;
	const FENCE_GATE_SPRUCE = 183;
	const FENCE_GATE_BIRCH = 184;
	const FENCE_GATE_JUNGLE = 185;
	const FENCE_GATE_DARK_OAK = 186;
	const FENCE_GATE_ACACIA = 187;
	const GRASS_PATH = 198;
	const PODZOL = 243;
	const BEETROOT_BLOCK = 244;
	const STONECUTTER = 245;
	const GLOWING_OBSIDIAN = 246;
	//Normal Item IDs
	const IRON_SHOVEL = 256; //
	const IRON_PICKAXE = 257; //
	const IRON_AXE = 258; //
	const FLINT_STEEL = 259; //
	const FLINT_AND_STEEL = 259; //
	const APPLE = 260; //
	const BOW = 261;
	const ARROW = 262;
	const COAL = 263; //
	const DIAMOND = 264; //
	const IRON_INGOT = 265; //
	const GOLD_INGOT = 266; //
	const IRON_SWORD = 267;
	const WOODEN_SWORD = 268; //
	const WOODEN_SHOVEL = 269; //
	const WOODEN_PICKAXE = 270; //
	const WOODEN_AXE = 271; //
	const STONE_SWORD = 272;
	const STONE_SHOVEL = 273;
	const STONE_PICKAXE = 274;
	const STONE_AXE = 275;
	const DIAMOND_SWORD = 276;
	const DIAMOND_SHOVEL = 277;
	const DIAMOND_PICKAXE = 278;
	const DIAMOND_AXE = 279;
	const STICK = 280; //
	const STICKS = 280;
	const BOWL = 281; //
	const MUSHROOM_STEW = 282;
	const GOLD_SWORD = 283;
	const GOLD_SHOVEL = 284;
	const GOLD_PICKAXE = 285;
	const GOLD_AXE = 286;
	const GOLDEN_SWORD = 283;
	const GOLDEN_SHOVEL = 284;
	const GOLDEN_PICKAXE = 285;
	const GOLDEN_AXE = 286;
	const STRING = 287;
	const FEATHER = 288; //
	const GUNPOWDER = 289;
	const WOODEN_HOE = 290;
	const STONE_HOE = 291;
	const IRON_HOE = 292; //
	const DIAMOND_HOE = 293;
	const GOLD_HOE = 294;
	const GOLDEN_HOE = 294;
	const SEEDS = 295;
	const WHEAT_SEEDS = 295;
	const WHEAT = 296;
	const BREAD = 297;
	const LEATHER_CAP = 298;
	const LEATHER_TUNIC = 299;
	const LEATHER_PANTS = 300;
	const LEATHER_BOOTS = 301;
	const CHAIN_HELMET = 302;
	const CHAIN_CHESTPLATE = 303;
	const CHAIN_LEGGINGS = 304;
	const CHAIN_BOOTS = 305;
	const IRON_HELMET = 306;
	const IRON_CHESTPLATE = 307;
	const IRON_LEGGINGS = 308;
	const IRON_BOOTS = 309;
	const DIAMOND_HELMET = 310;
	const DIAMOND_CHESTPLATE = 311;
	const DIAMOND_LEGGINGS = 312;
	const DIAMOND_BOOTS = 313;
	const GOLD_HELMET = 314;
	const GOLD_CHESTPLATE = 315;
	const GOLD_LEGGINGS = 316;
	const GOLD_BOOTS = 317;
	const FLINT = 318;
	const RAW_PORKCHOP = 319;
	const COOKED_PORKCHOP = 320;
	const PAINTING = 321;
	const GOLDEN_APPLE = 322;
	const SIGN = 323;
	const WOODEN_DOOR = 324;
	const BUCKET = 325;
	const MINECART = 328;
	const IRON_DOOR = 330;
	const REDSTONE = 331;
	const REDSTONE_DUST = 331;
	const SNOWBALL = 332;
	const LEATHER = 334;
	const BRICK = 336;
	const CLAY = 337;
	const SUGARCANE = 338;
	const SUGAR_CANE = 338;
	const SUGAR_CANES = 338;
	const PAPER = 339;
	const BOOK = 340;
	const SLIMEBALL = 341;
	const EGG = 344;
	const COMPASS = 345;
	const FISHING_ROD = 346;
	const CLOCK = 347;
	const GLOWSTONE_DUST = 348;
	const RAW_FISH = 349;
	const COOKED_FISH = 350;
	const DYE = 351;
	const BONE = 352;
	const SUGAR = 353;
	const CAKE = 354;
	const BED = 355;
	const COOKIE = 357;
	const SHEARS = 359;
	const MELON = 360;
	const MELON_SLICE = 360;
	const PUMPKIN_SEEDS = 361;
	const MELON_SEEDS = 362;
	const RAW_BEEF = 363;
	const STEAK = 364;
	const COOKED_BEEF = 364;
	const RAW_CHICKEN = 365;
	const COOKED_CHICKEN = 366;
	const GOLD_NUGGET = 371;
	const GOLDEN_NUGGET = 371;
	const SPAWN_EGG = 383;
	const EMERALD = 388;
	const FLOWER_POT = 390;
	const CARROT = 391;
	const CARROTS = 391;
	const POTATO = 392;
	const POTATOES = 392;
	const BAKED_POTATO = 393;
	const BAKED_POTATOES = 393;
	const PUMPKIN_PIE = 400;
	const NETHER_BRICK = 405;
	const QUARTZ = 406;
	const NETHER_QUARTZ = 406;
	const CAMERA = 456;
	const BEETROOT = 457;
	const BEETROOT_SEEDS = 458;
	const BEETROOT_SEED = 458;
	const BEETROOT_SOUP = 459;
'''

def UpperToPascal(name):
    slices = name.split("_")
    ret = ''
    for s in slices:
        ret += s[0].upper()
        ret += s[1:].lower()

    return ret

if __name__ == "__main__":
    toimport = re.sub("([\t\r;(//)])", "", toimport).split("\n")
    numconst = dict()
    constnum = dict()
    alias = dict()
    constdef = ''
    mapdef = ''
    for line in toimport:
        if len(line) < 5 or line[:5] != "const":
            continue

        p = line[6:].split(" = ")
        p = (UpperToPascal(p[0]), int(p[1]))
        constnum[p[0]] = p[1]
        if p[1] in numconst:
            alias[p[0]] = numconst[p[1]]
        else:
            numconst[p[1]] = p[0]

    print("""package types

// ID represents ID for Minecraft blocks/items.
type ID uint16

// String converts ID to string.
func (id ID) String() string {
    if name, ok := nameMap[id]; ok {
        return name
    }
    return "Unknown"
}

// Block tries to convert item ID to block ID. If fails, it panics.
func (id ID) Block() byte {
    if id >= 256 {
        panic(fmt.Sprintf("item ID %d(%s) overflows byte", uint16(id), id))
    }
    return byte(id)
}

    """)

    print("// item/block IDs\nconst (\n    %s ID = iota" % numconst.pop(0, None))
    last = 0
    for i in numconst:
        c = numconst[i]
        for ii in range((i - last - 1)):
            print("    _ // %d" % (last + ii + 1))
        print("    %s // %d" % (c, i))
        last = i
    print(")\n")
    print("// aliases\nconst (")
    for n in alias:
        a = alias[n]
        print("    %s = %s // %d" % (n, a, constnum[a]))
    print(")")

    print("var idMap = map[string]ID {")
    for i in numconst:
        c = numconst[i]
        print("    \"%s\": %s, // %d" % (c, c, i))

    print("\n    //aliases\n")
    for n in alias:
        i = constnum[alias[n]]
        print("    \"%s\": %s, // %d" % (n, n, i))
    print("}\n")

    print("var nameMap = map[ID]string {")
    for i in numconst:
        c = numconst[i]
        print("    %s: \"%s\", // %d" % (c, c, i))
    print("}\n")
    print("""
func StringID(name string) ID {
    if id, ok := idMap[name]; ok {
        return id
    }
    return 65535
}
""")
