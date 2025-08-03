import { CommandInteraction, GuildMember, MessageFlags, SlashCommandBuilder } from "discord.js";
import { APIEndpoints, Request } from "../../api";
import { dconfig } from "../../config";
import { CacheUtil } from "../../core/CacheUtil";
import { EmbedUtil } from "../../core/EmbedUtil";

export const data = new SlashCommandBuilder()
    .setName("unregister")
    .setDescription("To unregister your Discord account from your MC account")
    .addMentionableOption(option => option.setName("member").setDescription("Input member you want to unregister").setRequired(true));

export async function execute(interaction: CommandInteraction) {
    const member = interaction.options.get("member")?.member as GuildMember;
    const res = await Request.get(`${APIEndpoints.GET_REGISTERED_PLAYER}/${member.user.id}`);
    if (!res.success) {
        await interaction.reply({
            embeds: [EmbedUtil.create({
                type: "no",
                description: `<@${member.id}> is already unregistered.`,
            })],
            flags: MessageFlags.Ephemeral
        });
        return;
    }

    res.data.userid = "";

    await Request.post(APIEndpoints.UPDATE_PLAYER, res.data);

    member.roles.remove(CacheUtil.getRole(member.guild, dconfig.roles.registered));
    member.setNickname(null).catch(() => { });
    await interaction.reply({
        embeds: [EmbedUtil.create({
            type: "yes",
            description: `Successfully unregistered <@${member.id}>. Their nickname has been reset and their <#${dconfig.roles.registered}> role removed.`,
        })],
        flags: MessageFlags.Ephemeral
    });
}
