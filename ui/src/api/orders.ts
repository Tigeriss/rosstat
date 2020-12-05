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
    type: number;
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

//

export interface BigPalletModel {
    pallet_num: number;
    types: BigOrdersModel[];
}

export interface BigPalletBarcodeModel {
    success: boolean;
    type: number;
    error: string;
}

export interface BigPalletFinishRequestModel {
    pallet_num: number;
    barcodes: string[];
}

export interface BigPalletFinishResponseModel {
    success: boolean;
    error: string;
    last_pallet: boolean;
}

export async function getBigPallet(session: Session, id: number): Promise<BigPalletModel> {
    return await request(session, `orders/big/pallet/${id}`, {}, "GET");
}

export async function getBigPalletBarcode(session: Session, id: number, barcode: string): Promise<BigPalletBarcodeModel> {
    return await request(session, `orders/big/pallet/${id}/barcode/${barcode}`, {}, "GET");
}

export async function finishBigPallet(session: Session, id: number, req: BigPalletFinishRequestModel): Promise<BigPalletFinishResponseModel> {
    return await request(session, `orders/big/pallet/${id}/finish`, req);
}
