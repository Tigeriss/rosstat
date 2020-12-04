import {Session} from "../store/session";
import {request} from "./request";

export interface SubOrderModel {
    is_small: boolean;
    order_caption: string;
    amount_pallets: number;
    amount_boxes: number;
}

export interface OrdersModel {
    id: number;
    num: number;
    order_caption: string;
    customer: string;
    address: string;
    run: number;
    amount_pallets: number;
    amount_boxes: number;
    sub_orders: SubOrderModel[];
    opened: boolean;
}


export async function getOrdersToBuild(session: Session): Promise<OrdersModel[]> {
    return await request(session, "orders/build", {}, "GET");
}

export interface BigOrdersModel {
    form_name: string;
    total: number;
    built: number;
}

export async function getBigOrdersToBuild(session: Session, id: number): Promise<BigOrdersModel[]> {
    return await request(session, `orders/big/build/${id}`, {}, "GET");
}
