import { CategoryChannel, Guild, GuildTextBasedChannel, Role } from "discord.js";

export class CacheUtils {
    static getCategory(guild: Guild, categoryId: string): CategoryChannel {
        return guild.channels.cache.get(categoryId) as CategoryChannel;
    }

    static getChannel(guild: Guild, channelId: string): GuildTextBasedChannel {
        return guild.channels.cache.get(channelId) as GuildTextBasedChannel;
    }

    static getRole(guild: Guild, roleId: string): Role {
        return guild.roles.cache.get(roleId) as Role;
    }
}