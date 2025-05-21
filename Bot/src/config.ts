import dotenv from "dotenv";

dotenv.config();

const { DISCORD_TOKEN, JWT_SECRET } = process.env;

if (!DISCORD_TOKEN) {
    throw new Error("Missing environment variables");
}

export const token = DISCORD_TOKEN;
export const jwtSecret = JWT_SECRET;
export const config = {
    prefix: "=",
    clientId: "1268295659820159118",
    guildId: "1373222394046578731"
}
