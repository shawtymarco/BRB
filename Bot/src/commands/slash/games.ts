import { ChatInputCommandInteraction, GuildMember, MessageFlags, SlashCommandBuilder } from "discord.js";
import { gamesDB } from "../../core/GameCore";
import { EmbedUtil } from "../../core/EmbedUtil";
// in complete, no test yet
export const data = new SlashCommandBuilder()
    .setName("games")
    .setDescription("Show the games a user has played this season")
    .addUserOption(o => o.setName('user').setDescription('Target user'));

export async function execute(interaction: ChatInputCommandInteraction) {
    const user = (interaction.options.getUser('user') ?? (interaction.member as GuildMember).user);

    const active = [...gamesDB.data.values()].filter(g => (g as any)["memberIds"].includes(user.id));
    const lines = await Promise.all(active.map(async g => {
        return `#${(g as any).id.slice(0,6).toUpperCase()} }`;
    }));

    await interaction.reply({
        embeds: [EmbedUtil.create({
            type: 'yes',
            description: lines.length ? lines.join('\n\n') : 'No games found for this user in the current session.'
        })],
        flags: MessageFlags.Ephemeral
    });
}


