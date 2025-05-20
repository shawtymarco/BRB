import { Message } from "discord.js";

export const name = "cg";

export async function execute(message: Message, args: any) {
    await message.reply("WIP");
}