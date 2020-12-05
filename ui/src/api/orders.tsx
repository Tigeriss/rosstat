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


export async function getSmallOrdersToBuild(session: Session, id: number): Promise<BigOrdersModel[]> {
    return await request(session, `orders/small/build/${id}`, {}, "GET");
}

export async function finishOrders(session: Session, orderId: number, preparedBoxes: string[]): Promise<void> {
    return await request(session, `orders/small/build/${orderId}/finish`, {
        boxes: preparedBoxes
    }, "POST");
}
