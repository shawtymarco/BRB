import fs from "fs";
import https from "https";
import axios from "axios";
import { jwtSecret } from "./config";

const httpsAgent = new https.Agent({
    cert: fs.readFileSync('../certs/client.crt'),
    key: fs.readFileSync('../certs/client.key'),
    ca: fs.readFileSync('../certs/ca.pem'),
    rejectUnauthorized: true
});

export class Request {
    static async get(endpoint: string) {
        try {
            const res = await axios.get(`https://localhost:8080/${endpoint}`, {
                headers: {
                    'Authorization': `Bearer ${jwtSecret}`,
                },
                httpsAgent
            });

            console.log('GET response:', res.data);
        } catch (err: any) {
            console.error('Failed to send GET:', err.message);
          }
    }
}