import { ChatInputCommandInteraction, GuildMember, MessageFlags, SlashCommandBuilder } from "discord.js";
import { APIEndpoints, Request } from "../../api";
import { EmbedUtil } from "../../core/EmbedUtil";
import { Game } from "../../core/GameCore";

export const data = new SlashCommandBuilder()
    .setName("refresh")
    .setDescription("Refreshes your Discord profile with updated elo, username, and rank");

export async function execute(interaction: ChatInputCommandInteraction) {
    const member = interaction.member as GuildMember;
    
    await interaction.deferReply({ flags: MessageFlags.Ephemeral });

    try {
        const res = await Request.get(`${APIEndpoints.GET_REGISTERED_PLAYER}/${member.user.id}`);
        
        if (!res.data) {
            return interaction.editReply({
                embeds: [EmbedUtil.create({
                    type: "no",
                    description: "You are not registered! Use `/register` to link your Discord account with your Minecraft account first.",
                })]
            });
        }

        await Game.refreshMemberNickname(member);
        
        await Game.refreshMemberRank(member);

        return interaction.editReply({
            embeds: [EmbedUtil.create({
                type: "yes",
                description: `Successfully refreshed your profile!\n**Username:** ${res.data.Username}\n**ELO:** ${res.data.Statistics.ELO}`,
            })]
        });

    } catch (error) {
        console.error("Error refreshing user data:", error);
        return interaction.editReply({
            embeds: [EmbedUtil.create({
                type: "no",
                description: "Failed to refresh your profile. Please try again later or contact an administrator if the issue persists.",
            })]
        });
    }
}
