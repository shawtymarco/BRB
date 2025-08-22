import dotenv from "dotenv";

dotenv.config();

const { DISCORD_TOKEN, JWT_SECRET } = process.env;

if (!DISCORD_TOKEN) {
    throw new Error("Missing environment variables");
}

export const token = DISCORD_TOKEN;
export const jwtSecret = JWT_SECRET;
export const dconfig = {
    prefix: "=",
    api: "https://localhost:8080",
    clientId: "1268295659820159118",
    guildId: "967360687561138228",

    roles: {
        registered: "1258628748878676029",
        darkMatter: "1265242510490603622",
        topaz: "1265242500160032829",
        quartz: "1265242481205706753",
        aventurine: "1265242464567033887",
        obsidian: "1265241983379570772",
        amethyst: "1265241969647423499",
        opal: "1265241953625444384",
        crystal: "1265241933870272583",
        ruby: "1265241758061826108",
        sapphire: "1265241620501233664",
        emerald: "1265241483376853095",
        diamond: "1265241154471858186",
        platinum: "1265241055734005822",
        gold: "1265240921876725892",
        silver: "1265240836178837557",
        bronze: "1265240685531889665"
    },

    categories: {
        games: "1330480792777789500"
    },

    channels: {
        register: "1406955689561030797",
        touchAlerts: "1234129493535228046",
        touch2v2: "1339305887084318820",
        touch3v3: "1336343043283881984",
        allAlerts: "1234129493535228046",
        all3v3: "1234129207710318592",
        all4v4: "1234129355194765443",
        gameChat: "1234130482455908413",
        waitingRoom: "1234130106684018830",
        scoring: "1234130482455908413"
    },

    emojis: {
        yes: "<:yes:1360573400908824637>",
        no: "<:no:1360573506865332387>"
    }
}

export const dconfig_test = {
    prefix: "=",
    api: "https://localhost:8080",
    clientId: "1268295659820159118",
    guildId: "1373222394046578731",

    roles: {
        registered: "1400435527039062087"
    },

    categories: {
        games: "1381034727867416708"
    },

    channels: {
        register: "1375134986986065990",
        touchAlerts: "1375129254425268294",
        touch2v2: "1375128999260327976",
        touch3v3: "1375129021959901257",
        allAlerts: "1375129218849181767",
        all3v3: "1375128890846216202",
        all4v4: "1375128951244193792",
        gameChat: "1381376097161187398",
        waitingRoom: "1390492591643955210",
        scoring: "1390499656408105011"
    },

    emojis: {
        yes: "<:check_:1375131098392039536>",
        no: "<:No:1375131153576628344>"
    }
}
