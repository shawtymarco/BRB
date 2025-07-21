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
    guildId: "1373222394046578731",

    roles: {
        registered: "1375112083103813652"
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
