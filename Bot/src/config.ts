import dotenv from "dotenv";

dotenv.config();

const { DISCORD_TOKEN } = process.env;

if (!DISCORD_TOKEN) {
    throw new Error("Missing environment variables");
}

export const token = DISCORD_TOKEN;
export const config = {
    prefix: "=",
    clientId: "1268295659820159118",
    guildId: "1373222394046578731"
}
