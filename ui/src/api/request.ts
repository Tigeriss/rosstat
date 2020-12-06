import {Session} from "../store/session";

const API_URL = "/api/";

export async function request<T>(session: Session, path: string, data: unknown, method: string = "POST"): Promise<T> {

    const headers: Record<string, string> = {};
    headers["Content-Type"] = "application/json";

    if (session.currentUser != null) {
        headers["Authorization"] = `Bearer ${session.currentUser.token}`;
    }

    const request: RequestInit = {
        method: method,
        headers: headers,
    };

    if (method === "POST") {
        request.body = JSON.stringify(data);
    }

    const req = await fetch(`${API_URL + path}`, request);

    if (!req.ok) {
        throw new Error(`request ${method} ${API_URL + path} failed: ${req.statusText}`);
    }

    if (req.status === 204) {
        return {} as any;
    }

    return await req.json();
}
