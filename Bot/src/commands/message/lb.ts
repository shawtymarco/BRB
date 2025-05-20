import { Message } from "discord.js";

export const name = "lb";

export async function execute(message: Message, args: any) {
    await message.reply("WIP");
}