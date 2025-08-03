import { CommandInteraction, GuildMember, MessageFlags, SlashCommandBuilder } from "discord.js";
import { APIEndpoints, Request } from "../../api";
import { dconfig } from "../../config";
import { EmbedUtil } from "../../core/EmbedUtil";
import { CacheUtil } from "../../core/CacheUtil";

export const data = new SlashCommandBuilder()
    .setName("register")
    .setDescription("To link your Discord account with your MC account")
    .addStringOption(option => option.setName("code").setDescription("Input code shown after executing the command /link in-game").setMaxLength(4).setRequired(true));

export async function execute(interaction: CommandInteraction) {
    const member = interaction.member as GuildMember;
    const code = interaction.options.get("code");
    const res = await Request.post(APIEndpoints.VERIFY, { userId: member.user.id, code: code?.value });

    if (res.success) {
        member.roles.add(CacheUtil.getRole(member.guild, dconfig.roles.registered));
        member.setNickname(`${res.elo} 〣 ${res.username}`).catch(() => { });
    }

    await interaction.reply({
        embeds: [EmbedUtil.create({
            type: res.success ? "yes" : "no",
            description: res.message,
        })],
        flags: MessageFlags.Ephemeral
    });
}
