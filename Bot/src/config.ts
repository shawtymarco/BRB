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
    clientId: "1268295659820159118",
    guildId: "1373222394046578731",

    roles: {
        registered: "1375112083103813652"
    },

    channels: {
        register: "1375134986986065990"
    },

    emojis: {
        yes: "<:check_:1375131098392039536>",
        no: "<:No:1375131153576628344>"
    }
}
