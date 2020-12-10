import {Session} from "../store/session";
import {request} from "./request";

export interface ShipmentModel {
    id: number;
    num: number;
    order_caption: string;
    customer: string;
    address: string;
    run: number;
    amount_pallets: number;
    amount_boxes: number;
    collected: boolean;
    shipped: boolean;
}

export async function getShipmentReady(session: Session): Promise<ShipmentModel[]> {
    return await request(session, "shipment/ready", {}, "GET");
}

export interface ShipmentPalletModel {
    num: number;
    pallet_num: number;
    barcode: string;
    amount_boxes: number;
}

export interface ShipmentReportModel {
    order_caption: string;
    address: string;
    total_boxes: number;
    total_pallets: number;
    items: ShipmentReportItemModel[];
}

export interface ShipmentReportItemModel {
    num: number;
    name: string;
    run: number;
    amount_in_box: number;
    completed_boxes: number;
    amount_in_composed_box: number;
}

export async function getShipmentPallet(session: Session, id: number): Promise<ShipmentPalletModel[]> {
    return await request(session, `shipment/pallet/${id}`, {}, "GET");
}

export async function finishPalletShipment(session: Session, id: number): Promise<void> {
    return await request(session, `shipment/pallet/${id}/finish`, {});
}

export async function getShipmentReport(session: Session, id: number): Promise<ShipmentReportModel> {
    return await request(session, `shipment/pallet/${id}/report`, {}, "GET");
}
