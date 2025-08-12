import {
    ActionRowBuilder,
    ChannelType,
    Collection,
    EmbedBuilder,
    Guild,
    GuildMember,
    OverwriteResolvable,
    OverwriteType,
    PermissionFlagsBits,
    PrivateThreadChannel,
    StringSelectMenuBuilder,
    StringSelectMenuOptionBuilder,
    VoiceChannel
} from "discord.js";
import path from "path";
import { client } from "..";
import { APIEndpoints, Request } from "../api";
import { dconfig } from "../config";
import { CacheUtil } from "./CacheUtil";
import { DB } from "./DB";

export var gamesDB: DB<Game> = new DB(path.join(".", "db", "games.json"));

export class Game {
    team1Ids: string[] = [];
    team2Ids: string[] = [];
    teamPickingMessageId: string | null = null;
    step: number = 0;

    constructor(
        public id: string,
        private threadId: string,
        private lobbyVCId: string,
        private team1VCId: string,
        private team2VCId: string,
        public teamSize: number,
        private memberIds: string[],
        public captainIds: string[]
    ) {
        this.team1Ids.push(captainIds[0]);
        this.team2Ids.push(captainIds[1]);

        const interval = setInterval(async () => {
            try {
                const res = await Request.get(APIEndpoints.GET_GAMES_TO_TERMINATE);
                const gameData = res.ids[this.id];
                if (gameData) {
                    clearInterval(interval);
                    await this.terminateGame(gameData);
                }
            } catch (err) {
                console.error(`Failed to connect to ${dconfig.api}`);
            }
        }, 5000);
    }

    get guild(): Guild {
        return client.guilds.cache.get(dconfig.guildId) as Guild
    }

    thread(): PrivateThreadChannel {
        return this.guild.channels.cache.get(this.threadId) as PrivateThreadChannel;
    }

    lobbyVC(): VoiceChannel {
        return this.guild.channels.cache.get(this.lobbyVCId) as VoiceChannel;
    }

    team1VC(): VoiceChannel {
        return this.guild.channels.cache.get(this.team1VCId) as VoiceChannel;
    }

    team2VC(): VoiceChannel {
        return this.guild.channels.cache.get(this.team2VCId) as VoiceChannel;
    }

    async members(): Promise<Collection<string, GuildMember>> {
        const res = new Collection<string, GuildMember>();

        const fetchedMembers = await Promise.all(
            this.memberIds.map(id => this.guild.members.fetch(id).then(member => [id, member] as const))
        );

        for (const [id, member] of fetchedMembers) {
            res.set(id, member);
        }

        return res;
    }

    async team1(): Promise<GuildMember[]> {
        return await Promise.all(
            this.team1Ids.map(id => this.guild.members.fetch(id))
        );
    }

    async team2(): Promise<GuildMember[]> {
        return await Promise.all(
            this.team2Ids.map(id => this.guild.members.fetch(id))
        );
    }

    async captains(): Promise<GuildMember[]> {
        return await Promise.all(
            this.captainIds.map(id => this.guild.members.fetch(id))
        );
    }

    isTeam1Turn(): Boolean {
        return this.step % 2 === 0;
    }

    async sendIntroductionMessage() {
        const captains = await this.captains();

        this.thread().send({
            embeds: [{
                author: {
                    name: `Bedrock Ranked Bedwars | Game #${this.id}`,
                    icon_url: 'https://images-ext-1.discordapp.net/external/xPUGYxZAJDXj4ScgckfwI0SvwkRQDNQDTi2gF27kRNc/%3Fsize%3D4096/https/cdn.discordapp.com/avatars/1209943786252144690/6d272ead1117efac2cf674582fceff1f.png?format=webp&quality=lossless&width=801&height=801'
                },
                description: `**Matchmaking Type:** Captain
                    > 🎲 Random Captains have been chosen!
    
                    **Team 1 Captain:** <@${captains[0].id}>
                    **Team 2 Captain:** <@${captains[1].id}>`,
                color: 0xFFFFFF,
                footer: {
                    text: `brbw.net`
                },
                timestamp: new Date().toISOString()
            }]
        });
    }

    async updateCaptainPickingMessage() {
        const picks = [
            [1, 1],
            [1, 1, 1, 1],
            [1, 2, 2, 1]
        ][this.teamSize - 2][this.step];
        const thread = this.thread();
        const members = await this.members();
        const captains = await this.captains();
        const team1 = await this.team1();
        const team2 = await this.team2();

        const remainingMembers = members.filter(member => !team1.includes(member) && !team2.includes(member));
        if (remainingMembers.size == 1 && this.teamPickingMessageId != null) {
            (this.isTeam1Turn() ? this.team1Ids : this.team2Ids).push((remainingMembers.first() as GuildMember).id)

            const teamPickingMessage = await thread.messages.fetch(this.teamPickingMessageId);
            teamPickingMessage.edit({
                content: '',
                embeds: [{
                    author: {
                        name: `Bedrock Ranked Bedwars | Team Selection [Ended]`,
                        icon_url: 'https://images-ext-1.discordapp.net/external/xPUGYxZAJDXj4ScgckfwI0SvwkRQDNQDTi2gF27kRNc/%3Fsize%3D4096/https/cdn.discordapp.com/avatars/1209943786252144690/6d272ead1117efac2cf674582fceff1f.png?format=webp&quality=lossless&width=801&height=801'
                    },
                    fields: [
                        {
                            name: "Team 1",
                            value: `${this.team1Ids.map(id => `<@${id}>`).join("\n")}`,
                            inline: true,
                        },
                        {
                            name: "Team 2",
                            value: `${this.team2Ids.map(id => `<@${id}>`).join("\n")}`,
                            inline: true,
                        },
                    ],
                    description: `The captains finished picking the teams.`,
                    color: 0xFFFFFF,
                    footer: {
                        text: `brbw.net`
                    },
                    timestamp: new Date().toISOString()
                }],
                components: []
            });
            await this.concludeGameMaking();
            gamesDB.save();
            return;
        }

        const msgObj = {
            content: `<@${this.isTeam1Turn() ? captains[0].id : captains[1].id}>`,
            embeds: [{
                author: {
                    name: `Bedrock Ranked Bedwars | Team Selection`,
                    icon_url: 'https://images-ext-1.discordapp.net/external/xPUGYxZAJDXj4ScgckfwI0SvwkRQDNQDTi2gF27kRNc/%3Fsize%3D4096/https/cdn.discordapp.com/avatars/1209943786252144690/6d272ead1117efac2cf674582fceff1f.png?format=webp&quality=lossless&width=801&height=801'
                },
                fields: [
                    {
                        name: "Team 1",
                        value: `${team1.map(member => `<@${member.id}>`).join("\n")}`,
                        inline: true,
                    },
                    {
                        name: "Team 2",
                        value: `${team2.map(member => `<@${member.id}>`).join("\n")}`,
                        inline: true,
                    },
                ],
                description: `Hey <@${this.isTeam1Turn() ? captains[0].id : captains[1].id}>!\nIt's your turn to pick ${picks === 1 ? '**ONE** teammate' : '**TWO** teammates'}.`,
                color: 0xFFFFFF,
                footer: {
                    text: `brbw.net`
                },
                timestamp: new Date().toISOString()
            }],
            components: [new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(
                new StringSelectMenuBuilder()
                    .setCustomId("pick_teammates")
                    .setPlaceholder("Pick a teammate")
                    .setMinValues(1)
                    .setMaxValues(picks)
                    .addOptions(...remainingMembers.map(member => new StringSelectMenuOptionBuilder()
                        .setLabel(`${member.nickname} [${member.user.username}]`)
                        .setValue(member.id))
                    ))]
        };

        if (this.teamPickingMessageId != null) {
            const teamPickingMessage = await thread.messages.fetch(this.teamPickingMessageId);
            teamPickingMessage.edit(msgObj);
        } else {
            const msg = await thread.send(msgObj);
            this.teamPickingMessageId = msg.id;
        }
        gamesDB.save();
    }

    async concludeGameMaking() {
        const thread = this.thread();
        const members = await this.members();
        const t1 = await this.team1();
        const t2 = await this.team2();

        await this.createTeamVCs();

        members.forEach(member => {
            member.voice.setChannel(t1.includes(member) ? this.team1VC() : this.team2VC())
        })

        try {
            await Request.post(APIEndpoints.GAME_CONNECT_USERS, {
                Users: [...t1.map(m => m.id), ...t2.map(m => m.id)],
                code: this.id
            });
        } catch (err) {
            console.log(t1.map(m => m.id));
            console.log(t2.map(m => m.id));
            console.log([...t1.map(m => m.id), ...t2.map(m => m.id)])
        }

        thread.send({
            content: `${members.map(m => `<@${m.id}>`).join(" ")}\n## Team making is DONE. Join as soon as possible to start playing!`,
            embeds: [{
                author: {
                    name: `Bedrock Ranked Bedwars | How to join`,
                    icon_url: 'https://images-ext-1.discordapp.net/external/xPUGYxZAJDXj4ScgckfwI0SvwkRQDNQDTi2gF27kRNc/%3Fsize%3D4096/https/cdn.discordapp.com/avatars/1209943786252144690/6d272ead1117efac2cf674582fceff1f.png?format=webp&quality=lossless&width=801&height=801'
                },
                description: `
                - Join the server (brbw.net) and you will automatically be warped to the game.\n**OR**\n- In the lobby, type the command \`/join\` to manually join your game.

                Once there, you can execute \`/warp\` to warp players from the lobby into your game!
                `,
                color: 0xFFFFFF,
                footer: {
                    text: `brbw.net`
                },
                timestamp: new Date().toISOString()
            }]
        });

        await gamesDB.save();
    }

    async createTeamVCs() {
        const t1 = await this.team1();
        const t2 = await this.team2();

        const team1Permissions: OverwriteResolvable[] = [
            {
                id: this.guild.id,
                deny: [PermissionFlagsBits.ViewChannel, PermissionFlagsBits.Connect]
            },
        ];

        const team2Permissions: OverwriteResolvable[] = [
            {
                id: this.guild.id,
                deny: [PermissionFlagsBits.ViewChannel, PermissionFlagsBits.Connect]
            },
        ];

        t1.forEach(member => {
            team1Permissions.push({
                id: member.id,
                type: OverwriteType.Member,
                allow: [PermissionFlagsBits.ViewChannel, PermissionFlagsBits.Connect]
            });

            team2Permissions.push({
                id: member.id,
                type: OverwriteType.Member,
                allow: [PermissionFlagsBits.ViewChannel],
                deny: [PermissionFlagsBits.Connect]
            });
        });

        t2.forEach(member => {
            team2Permissions.push({
                id: member.id,
                type: OverwriteType.Member,
                allow: [PermissionFlagsBits.ViewChannel, PermissionFlagsBits.Connect]
            });

            team1Permissions.push({
                id: member.id,
                type: OverwriteType.Member,
                allow: [PermissionFlagsBits.ViewChannel],
                deny: [PermissionFlagsBits.Connect]
            })
        });

        const team1VC = await this.guild.channels.create({
            name: `#${this.id.slice(0, 4).toUpperCase()} | Team 1`,
            type: ChannelType.GuildVoice,
            parent: dconfig.categories.games,
            permissionOverwrites: team1Permissions
        });

        const team2VC = await this.guild.channels.create({
            name: `#${this.id.slice(0, 4).toUpperCase()} | Team 2`,
            type: ChannelType.GuildVoice,
            parent: dconfig.categories.games,
            permissionOverwrites: team1Permissions
        });

        this.team1VCId = team1VC.id;
        this.team2VCId = team2VC.id;

        await gamesDB.save();
    }

    async terminateGame(data: any) {
        const waitingRoom = CacheUtil.getChannel(this.guild, dconfig.channels.waitingRoom);

        if (this.thread()) {
            this.thread().setLocked(true);
            this.thread().setArchived(true);
            for (const [, member] of await this.thread().members.fetch()) {
                this.thread().members.remove(member.id);
            }
        }

        if (this.lobbyVC()) {
            this.lobbyVC().members.forEach(member => {
                member.voice.setChannel(waitingRoom as VoiceChannel);
            });
        }

        if (this.team1VC()) {
            this.team1VC().members.forEach(member => {
                member.voice.setChannel(waitingRoom as VoiceChannel);
            });
            this.team2VC().members.forEach(member => {
                member.voice.setChannel(waitingRoom as VoiceChannel);
            });
        }

        setTimeout(async () => {
            const lobbyVC = this.lobbyVC && this.lobbyVC();
            const team1VC = this.team1VC && this.team1VC();
            const team2VC = this.team2VC && this.team2VC();

            if (lobbyVC) await lobbyVC.delete();
            if (team1VC) await team1VC.delete();
            if (team2VC) await team2VC.delete();
        }, 3000);

        let embed;
        if (data) {
            embed = new EmbedBuilder()
                .setColor(0x2f3136)
                .setTitle(`#${this.id.slice(0, 4).toUpperCase()} - Bedrock Ranked Bedwars`)
                .addFields(
                    { name: "Game:", value: `#${this.id}`, inline: true },
                    { name: "Duration:", value: Game.formatDuration(data.Duration), inline: true },
                    { name: "MVP(s):", value: data.MVPs.map(m => `<@${m}>`).join(" "), inline: false },
                    { name: "Winning Team", value: Game.formatTeam(data.WinningTeam), inline: false },
                    { name: "Losing Team", value: Game.formatTeam(data.LosingTeam), inline: false },
                );
        } else {
            embed = new EmbedBuilder()
                .setColor(0x2f3136)
                .setTitle(`#${this.id.slice(0, 4).toUpperCase()} - Bedrock Ranked Bedwars`)
                .addFields(
                    { name: "Game:", value: `#${this.id}`, inline: true },
                    { name: "Status:", value: "**Voided**", inline: true },
                    { name: "Queue", value: this.memberIds.map(id => `<@${id}>`).join("\n"), inline: false },
                );
        }

        CacheUtil.getChannel(this.guild, dconfig.channels.scoring).send({
            content: this.memberIds.map(id => `<@${id}>`).join(""),
            embeds: [embed],
        });

        gamesDB.remove(this.threadId);
        await gamesDB.save();
    }

    static formatDuration(seconds: number): string {
        const h = Math.floor(seconds / 3600);
        const m = Math.floor((seconds % 3600) / 60);
        const s = Math.floor(seconds % 60);

        let parts: string[] = [];
        if (h > 0) parts.push(`${h}hr`);
        if (m > 0 || h > 0) parts.push(`${m}min`);
        parts.push(`${s}sec`);

        return parts.join(' ');
    }


    static formatTeam(team: Record<string, number[]>): string {
        return Object.entries(team).map(([id, [oldElo, newElo]]) => {
            const change = newElo - oldElo;
            const changeStr = change >= 0 ? `(+${change})` : `(${change})`;
            return `<@${id}> \`${changeStr}\` \`[${oldElo} ➝ ${newElo}]\``;
        }).join("\n");
    }

    static async refreshMemberNickname(member: GuildMember) {
        const res = await Request.get(`${APIEndpoints.GET_REGISTERED_PLAYER}/${member.user.id}`);
        member.setNickname(`${res.data.Statistics.ELO} 〣 ${res.data.Username}`).catch(() => { });
    }
}