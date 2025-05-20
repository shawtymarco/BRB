import { Message } from "discord.js";

export const name = "ping";

export async function execute(message: Message, args: any) {
    await message.reply("Pong!");
}