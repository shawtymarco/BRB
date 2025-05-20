import { GuildTextBasedChannel, Message, TextBasedChannel } from "discord.js";

export const name = "link";

export async function execute(message: Message<true>, args: any) {
    message.channel.send("jo")
    await message.reply("WIP");
}