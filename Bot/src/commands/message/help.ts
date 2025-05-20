import { Message } from "discord.js";

export const name = "help";

export async function execute(message: Message, args: any) {
    await message.reply("WIP");
}