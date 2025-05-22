import { GuildTextBasedChannel, Message, TextBasedChannel } from "discord.js";
import { APIEndpoints, Request } from "../../api";
import { dconfig } from "../../config";
import { EmbedUtil } from "../../core/EmbedUtil";

export const name = "register";

export async function execute(message: Message<true>, args: string[]) {
    if (message.channel.id !== dconfig.channels.register) {
        return;
    }

    const code = args[0];
    const res = await Request.post(APIEndpoints.VERIFY, { userId: message.author.id, code: code });

    if (res.success) {
        const role = message.guild.roles.cache.get(dconfig.roles.registered);
        if (role) message.member?.roles.add(role);

        console.log(res.elo, res.username);
        message.member?.setNickname(`${res.elo} 〣 ${res.username}`).catch(() => { });
    }

    const resMsg = await message.reply({
        embeds: [EmbedUtil.create({
            type: res.success ? "yes" : "no",
            description: res.message,
        })]
    });

    setTimeout(() => {
        message.delete().catch(() => { });
        resMsg.delete().catch(() => {});
    }, 5*1000);
}