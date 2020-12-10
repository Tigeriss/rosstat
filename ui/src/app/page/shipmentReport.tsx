import React, {useEffect, useState} from "react";
import {useSession} from "../app";
import {useParams} from "react-router-dom";
import {ShipmentReportItemModel, ShipmentReportModel} from "../../api/shipment";

function renderItem(item: ShipmentReportItemModel) {
    return <tr key={item.num} className="full-borders">
        <td>{item.num}</td>
        <td>{item.name}</td>
        <td>{item.run}</td>
        <td>{item.amount_in_box}</td>
        <td>{item.completed_boxes}</td>
        <td>{item.completed_boxes * item.amount_in_box}</td>
        <td>{item.amount_in_composed_box}</td>
    </tr>;
}

export const ShipmentReportPage = () => {
    const {id} = useParams<{ id: string }>();
    const session = useSession();
    const [report, setReport] = useState<ShipmentReportModel>({
        address: "",
        order_caption: "",
        total_boxes: 0,
        total_pallets: 0,
        items: []
    });

    useEffect(() => {
        (async () => {
            const res = await session.findShipmentReport(parseInt(id));
            setReport(res);
        })();
    }, [session, id]);

    return <div style={{padding: 30}}>
        <h1 style={{textAlign: "center"}}>Отчет по заказу № {report.order_caption}</h1>
        <h3>Адрес: {report.address}</h3>
        <table className="print">
            <thead>
            <tr>
                <th>Номер позиции</th>
                <th>Наименование</th>
                <th>Тираж</th>
                <th>В 1 коробке</th>
                <th>Всего целых коробок</th>
                <th>Кол-во в целых коробах</th>
                <th>Кол-во в сборных коробах</th>
            </tr>
            </thead>
            <tbody>
            {report.items.map(renderItem)}
            </tbody>
            <tfoot>
            <tr>
                <th colSpan={3} className="no-borders" />
                <th>Итого коробов</th>
                <th>{report.total_boxes}</th>
                <th colSpan={2} className="no-borders" />
            </tr>
            <tr>
                <th colSpan={3} className="no-borders"/>
                <th>Итого паллет</th>
                <th>{report.total_pallets}</th>
                <th colSpan={2} className="no-borders" />
            </tr>
            </tfoot>
        </table>
        <div>Дата: __________________________</div>
        <div style={{float: "right"}}>Подпись: __________________________</div>
        <div style={{clear: "both", height: 50}} />
        <div style={{float: "right", marginRight: 300}}>Место для печати</div>
    </div>;
}
