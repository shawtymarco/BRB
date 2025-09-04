import { ChatInputCommandInteraction, MessageFlags, SlashCommandBuilder } from "discord.js";
import { gamesDB, Game as GameModel } from "../../core/GameCore";
import { EmbedUtil } from "../../core/EmbedUtil";

export const data = new SlashCommandBuilder()
    .setName("game")
    .setDescription("Show info about a specified game")
    .addStringOption(o => o.setName('id').setDescription('Game ID').setRequired(true));

export async function execute(interaction: ChatInputCommandInteraction) {
    const id = interaction.options.get('id')?.value as string;

    const game = [...gamesDB.data.values()].find(g => (g as any)["id"] === id || (g as any)["id"].toString().startsWith(id));
    if (!game) {
        await interaction.reply({ embeds: [EmbedUtil.create({ type: 'no', description: 'Game not found.' })], flags: MessageFlags.Ephemeral });
        return;
    }

    const anyGame = game as any;
    const hasWinners = anyGame?.WinningTeam && Object.keys(anyGame.WinningTeam).length > 0;
    const hasLosers = anyGame?.LosingTeam && Object.keys(anyGame.LosingTeam).length > 0;
    const hasMVPs = Array.isArray(anyGame?.MVPs) && anyGame.MVPs.length > 0;

    if (hasWinners || hasLosers || hasMVPs) {
        await interaction.reply({
            embeds: [
                {
                    color: 0x2f3136,
                    title: `Game ${anyGame.id.slice(0,4).toUpperCase()} (Scored)`,
                    fields: [
                        { name: 'Game', value: `${anyGame.id.slice(0,4).toUpperCase()}`, inline: true },
                        ...(typeof anyGame?.Duration === 'number' ? [{ name: 'Duration', value: GameModel.formatDuration(anyGame.Duration), inline: true }] : []),
                        ...(hasMVPs ? [{ name: 'MVP(s)', value: anyGame.MVPs.map((m: string) => `<@${m}>`).join(' ') || '—', inline: false }] : []),
                        ...(hasWinners ? [{ name: 'Winners', value: GameModel.formatTeam(anyGame.WinningTeam) || '—', inline: false }] : []),
                        ...(hasLosers ? [{ name: 'Losers', value: GameModel.formatTeam(anyGame.LosingTeam) || '—', inline: false }] : []),
                    ]
                }
            ]
        });
        return;
    }

    const team1 = await game.team1();
    const team2 = await game.team2();

    const memberIds: string[] = (game as any).memberIds ?? [];
    const isPicking = (game.team1Ids.length + game.team2Ids.length) < memberIds.length;
    const phase = isPicking ? 'Team Picking' : 'Playing';

    await interaction.reply({
        embeds: [
            {
                color: 0x2f3136,
                title: `Game ${(game as any).id.slice(0,4).toUpperCase()}`,
                fields: [
                    { name: 'Game', value: `${(game as any).id.slice(0,4).toUpperCase()}`, inline: true },
                    { name: 'Team Size', value: `${game.teamSize}v${game.teamSize}`, inline: true },
                    { name: 'Phase', value: phase, inline: true },
                    ...(isPicking
                        ? [
                            { name: 'Queue', value: memberIds.map(id => `<@${id}>`).join('\n') || '—', inline: false },
                          ]
                        : [
                            { name: 'Team 1', value: team1.map(m => `<@${m.id}>`).join('\n') || '—', inline: true },
                            { name: 'Team 2', value: team2.map(m => `<@${m.id}>`).join('\n') || '—', inline: true },
                          ]
                    )
                ]
            }
        ]
    });
}


