import { ActionRowBuilder, ButtonBuilder, ButtonStyle, ChatInputCommandInteraction, SlashCommandBuilder } from "discord.js";
import { APIEndpoints, Request } from "../../api";
import { EmbedUtil } from "../../core/EmbedUtil";

const MODES = [
    "Deaths",
    "Kills",
    "Wins",
    "Losses",
    "WLR",
    "KDR",
    "Winstreak",
    "ELO",
    "Beds",
    "MVPs",
] as const;

type Mode = typeof MODES[number];

type Player = {
    Username: string;
    UserId: string;
    Statistics: { ELO: number };
    Games: {
        BedWars: { Wins: number; Losses: number; WinStreak: number; Kills: number; FinalKills: number; BedsBroken: number; Deaths: number; MVPCount: number };
    };
};

export const data = new SlashCommandBuilder()
    .setName("lb")
    .setDescription("Shows the leaderboard for a stat. Defaults to ELO.")
    .addStringOption(opt =>
        opt
            .setName("mode")
            .setDescription("Deaths, Kills, Wins, Losses, WLR, KDR, Winstreak, ELO, Beds, MVPs")
            .addChoices(...MODES.map(m => ({ name: m, value: m })))
    );

export async function execute(interaction: ChatInputCommandInteraction) {
    const selected = (interaction.options.getString("mode") as Mode | null) ?? "ELO";

    await interaction.deferReply();

    const res = await Request.get(APIEndpoints.GET_PLAYER);
    const players: Player[] = res.data ?? [];

    function statFor(p: Player, mode: Mode): number {
        const bw = p.Games.BedWars;
        switch (mode) {
            case "ELO": return p.Statistics.ELO ?? 0;
            case "Kills": return bw.Kills + bw.FinalKills;
            case "Deaths": return bw.Deaths;
            case "Wins": return bw.Wins;
            case "Losses": return bw.Losses;
            case "Beds": return bw.BedsBroken;
            case "MVPs": return bw.MVPCount ?? 0;
            case "Winstreak": return bw.WinStreak;
            case "KDR": return bw.Deaths === 0 ? (bw.Kills + bw.FinalKills) : (bw.Kills + bw.FinalKills) / Math.max(1, bw.Deaths);
            case "WLR": return bw.Losses === 0 ? bw.Wins : bw.Wins / Math.max(1, bw.Losses);
            default: return 0;
        }
    }

    const computed = players
        .filter(p => p.Username)
        .map(p => ({ name: p.Username, value: statFor(p, selected) }))
        .filter(x => x.value > 0)
        .sort((a, b) => b.value - a.value);

    const pageSize = 10;
    const totalPages = Math.max(1, Math.ceil(computed.length / pageSize));
    let page = 1;

    const medals = [":third_place:", ":second_place:", ":first_place:"];

    const buildDescription = () => {
        const start = (page - 1) * pageSize;
        const slice = computed.slice(start, start + pageSize);
        return slice
            .map((entry, idx) => {
                const rank = start + idx + 1;
                const prefix = rank <= 3 ? medals[3 - rank] : `#${rank}`;
                const shown = ["KDR", "WLR"].includes(selected) ? entry.value.toFixed(2) : Math.round(entry.value).toString();
                return `${prefix} ${entry.name} -> ${shown}`;
            })
            .join("\n")
            .concat(`\n\nPage ${page}/${totalPages} | Total Players: ${computed.length}`);
    };

    const embed = EmbedUtil.create({
        title: `Top ${selected} Leaderboard`,
        description: buildDescription()
    });

    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(
        new ButtonBuilder().setCustomId("lb_prev").setEmoji("◀️").setStyle(ButtonStyle.Secondary),
        new ButtonBuilder().setCustomId("lb_next").setEmoji("▶️").setStyle(ButtonStyle.Secondary),
    );

    const msg = await interaction.followUp({ embeds: [embed], components: totalPages > 1 ? [row] : [] });

    if (totalPages <= 1) return;

    const collector = msg.createMessageComponentCollector({ time: 60_000, filter: i => i.user.id === interaction.user.id });
    collector.on("collect", async i => {
        if (i.customId === "lb_prev") page = page <= 1 ? totalPages : page - 1;
        if (i.customId === "lb_next") page = page >= totalPages ? 1 : page + 1;
        await i.update({
            embeds: [EmbedUtil.create({ title: `Top ${selected} Leaderboard`, description: buildDescription() })]
        });
    });
}
