import { promises as fs } from "fs";
import { Collection } from "discord.js";

export class DB<T> {
    private filePath: string;
    public data: Collection<string, T> = new Collection();

    constructor(filePath: string) {
        this.filePath = filePath;
    }

    async load(mapper?: (raw: any) => T): Promise<void> {
        try {
            const raw = await fs.readFile(this.filePath, "utf8");
            const parsed: Record<string, unknown> = JSON.parse(raw);

            this.data = new Collection(
                Object.entries(parsed).map(([key, value]) => {
                    return [key, mapper ? mapper(value) : (value as T)];
                })
            );

            console.log(`[DB] Loaded ${this.data.size} entries from ${this.filePath}`);
        } catch (err: any) {
            if (err.code === "ENOENT") {
                console.log(`[DB] No existing file found (${this.filePath}), starting fresh.`);
                this.data = new Collection();
            } else {
                throw err;
            }
        }
    }

    async save(): Promise<void> {
        const obj: Record<string, T> = {};
        this.data.forEach((value, key) => {
            obj[key] = value;
        });

        await fs.writeFile(this.filePath, JSON.stringify(obj, null, 2));
        console.log(`[DB] Saved ${this.data.size} entries to ${this.filePath}`);
    }

    async add(key: string, entry: T): Promise<void> {
        this.data.set(key, entry);
        await this.save();
    }

    async update(key: string, updater: (item: T) => void): Promise<boolean> {
        const item = this.data.get(key);
        if (item) {
            updater(item);
            await this.save();
            return true;
        }
        return false;
    }

    async remove(key: string): Promise<boolean> {
        const removed = this.data.delete(key);
        if (removed) await this.save();
        return removed;
    }
}
