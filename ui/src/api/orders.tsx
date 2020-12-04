import {Session} from "../store/session";
import {request} from "./request";

export interface Orders {

}

export async function getOrders(session: Session): Promise<Orders> {
    return await request(session, "orders", {}, "GET");
}
