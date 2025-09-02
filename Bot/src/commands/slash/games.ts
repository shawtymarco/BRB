import { ChatInputCommandInteraction, GuildMember, MessageFlags, SlashCommandBuilder } from "discord.js";
import { gamesDB } from "../../core/GameCore";
import { EmbedUtil } from "../../core/EmbedUtil";
export const data = new SlashCommandBuilder()
    .setName("games")
    .setDescription("Show the games a user has played this season")
    .addUserOption(o => o.setName('user').setDescription('Target user'));

export async function execute(interaction: ChatInputCommandInteraction) {
    const user = (interaction.options.getUser('user') ?? (interaction.member as GuildMember).user);

    const relevant = [...gamesDB.data.values()].filter(g => {
        const anyGame = g as any;
        const memberIds: string[] = anyGame.memberIds ?? [];
        const hasScoreData = !!(anyGame?.WinningTeam || anyGame?.LosingTeam || anyGame?.MVPs);
        const inScored = (anyGame?.WinningTeam && user.id in anyGame.WinningTeam) || (anyGame?.LosingTeam && user.id in anyGame.LosingTeam) || (Array.isArray(anyGame?.MVPs) && anyGame.MVPs.includes(user.id));
        return memberIds.includes(user.id) || inScored || hasScoreData;
    });

    const lines = await Promise.all(relevant.map(async g => {
        const anyGame = g as any;
        const idStr = `#${(anyGame.id ?? '').slice(0,6).toUpperCase()}`;

        const hasWinners = anyGame?.WinningTeam && Object.keys(anyGame.WinningTeam).length > 0;
        const hasLosers = anyGame?.LosingTeam && Object.keys(anyGame.LosingTeam).length > 0;
        const hasMVPs = Array.isArray(anyGame?.MVPs) && anyGame.MVPs.length > 0;

        const won = !!(hasWinners && (user.id in anyGame.WinningTeam));
        const lost = !!(hasLosers && (user.id in anyGame.LosingTeam));
        const mvp = !!(hasMVPs && anyGame.MVPs.includes(user.id));

        const memberIds: string[] = anyGame.memberIds ?? [];
        const isPending = !won && !lost && !mvp && memberIds.includes(user.id);

        let icons: string[] = [];
        if (won) icons.push('🏆');
        if (mvp) icons.push('🔥');
        if (lost) icons.push('🔴');
        if (isPending) icons.push('🟡');
        if (icons.length === 0) icons.push('⚪');

        return `${icons.join('')} ${idStr} `;
    }));

    await interaction.reply({
        embeds: [EmbedUtil.create({
            type: 'yes',
            description: lines.length ? lines.join('\n\n') : 'No games found for this user in the current session.'
        })],
        flags: MessageFlags.Ephemeral
    });
}


