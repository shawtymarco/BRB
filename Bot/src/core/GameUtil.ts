import { GuildMember } from "discord.js";
import { APIEndpoints, Request } from "../api";

export class GameUtil {
    static async refreshMemberNickname(member: GuildMember) {
        const res = await Request.get(`${APIEndpoints.GET_REGISTERED_PLAYER}/${member.user.id}`);
        member.setNickname(`${res.data.Statistics.ELO} 〣 ${res.data.Username}`).catch(() => { });
    }
}