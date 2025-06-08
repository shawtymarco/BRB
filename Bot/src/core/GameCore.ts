import { ActionRowBuilder, BaseGuildVoiceChannel, Collection, Guild, GuildMember, Message, PrivateThreadChannel, PublicThreadChannel, StringSelectMenuBuilder, StringSelectMenuOptionBuilder } from "discord.js";
import { APIEndpoints, Request } from "../api";
import { DB } from "./DB";
import path from "path";
import { client } from "..";
import { dconfig } from "../config";
import { EmbedUtil } from "./EmbedUtil";

export var gamesDB: DB<Game> = new DB(path.join(".", "db", "games.json"));

export class Game {
    team1Ids: string[] = [];
    team2Ids: string[] = [];
    teamPickingMessageId: string | null = null;
    step: number = 0;

    constructor(
        public id: string,
        private threadId: string,
        private vcId: string,
        public teamSize: number,
        private memberIds: string[],
        public captainIds: string[]
    ) {
        this.team1Ids.push(captainIds[0]);
        this.team2Ids.push(captainIds[1]);
    }

    get guild(): Guild {
        return client.guilds.cache.get(dconfig.guildId) as Guild
    }

    thread(): PrivateThreadChannel {
        return this.guild.channels.cache.get(this.threadId) as PrivateThreadChannel;
    }

    vc(): BaseGuildVoiceChannel {
        return this.guild.channels.cache.get(this.vcId) as BaseGuildVoiceChannel;
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
                    text: `eliagic.club | <t:${Math.floor(Date.now() / 1000)}>`
                },
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

        if (remainingMembers.size == 0 && this.teamPickingMessageId != null) {
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
                            value: `${team1.map(member => `<@${member.id}>`).join("\n")}`,
                            inline: true,
                        },
                        {
                            name: "Team 2",
                            value: `${team2.map(member => `<@${member.id}>`).join("\n")}`,
                            inline: true,
                        },
                    ],
                    description: `The captains finished picking the teams.`,
                    color: 0xFFFFFF,
                    footer: {
                        text: `eliagic.club`
                    },
                }],
                components: []
            });
            await this.concludeGameMaking();
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
                    text: `eliagic.club`
                },
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
    }

    async concludeGameMaking() {
        const thread = this.thread();
        const members = await this.members();
        const t1 = await this.team1();
        const t2 = await this.team2();

        await Request.post(APIEndpoints.GAME_CONNECT_USERS, {
            Users: t1.map(m => m.id).push(... t2.map(m => m.id)),
            code: this.id
        });

        thread.send({
            content: `${members.map(m => `<@${m.id}>`).join(" ")}\n## Team making is DONE. Join as soon as possible to start playing!`,
            embeds: [{
                author: {
                    name: `Bedrock Ranked Bedwars | How to join`,
                    icon_url: 'https://images-ext-1.discordapp.net/external/xPUGYxZAJDXj4ScgckfwI0SvwkRQDNQDTi2gF27kRNc/%3Fsize%3D4096/https/cdn.discordapp.com/avatars/1209943786252144690/6d272ead1117efac2cf674582fceff1f.png?format=webp&quality=lossless&width=801&height=801'
                },
                description: `
                - Log into the server (eliagic.club) and you will automatically be teleported to the game.\n**OR**\n- In the lobby, type the command \`/join\` to manually join your game.

                Once there, wait up for the other players and you can execute \`/warp\` to teleport your game players from the lobby right into your game!
                `,
                color: 0xFFFFFF,
                footer: {
                    text: `eliagic.club`
                },
            }]
        })
    }

    static async refreshMemberNickname(member: GuildMember) {
        const res = await Request.get(`${APIEndpoints.GET_REGISTERED_PLAYER}/${member.user.id}`);
        member.setNickname(`${res.data.Statistics.ELO} 〣 ${res.data.Username}`).catch(() => { });
    }
}