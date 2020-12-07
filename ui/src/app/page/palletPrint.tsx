import React, {useEffect, useState} from "react";
import {useSession} from "../app";
import {useParams} from "react-router-dom";
import {PrintPalletModel, PrintPalletRegisterModel} from "../../api/orders";
import {Barcode} from "../component/barcode";

const stylesLabel: React.CSSProperties = {
    display: "inline-block",
    margin: "20px",
    padding: "30px 70px",
    border: "10px solid black",
    borderRadius: "50px",
};

const hStyles: React.CSSProperties = {
    fontSize: "2.8em"
}

function renderRow(reg: PrintPalletRegisterModel) {
    return <tr key={reg.num_pp}>
        <td>{reg.num_pp}</td>
        <td>{reg.position}</td>
        <td>{reg.amount}</td>
        <td>{reg.boxes}</td>
    </tr>;
}

export function PalletPrint() {
    const {id, num} = useParams<{ id: string, num: string }>();
    const [print, setPrint] = useState<Partial<PrintPalletModel>>({});
    const session = useSession();

    useEffect(() => {
        (async () => {
            const pallet = await session.findPallet(parseInt(id), parseInt(num));
            setPrint(pallet);
        })();
    }, [session, id, num]);

    return <div>
        <div style={stylesLabel} className="print-page">
            <h1 style={hStyles}>{print.order_caption}</h1>
            <h2 style={hStyles}>{print.address}</h2>
            <h2 style={hStyles}>{print.provider}</h2>
            <h2 style={hStyles}>{print.contract_number}</h2>
            <Barcode value={print.barcode ?? ""}/>
            <h2 style={hStyles}>Паллет №{num}</h2>
        </div>
        <div style={{padding: "30px"}}  className="print-page">
            <h1 style={{textAlign: "center"}}>Реестр паллеты №{num}</h1>
            <h1>Задание {print.order_caption}</h1>
            <table className="print">
                <thead>
                <tr>
                    <th>№п/п</th>
                    <th>Позиция</th>
                    <th>Штук</th>
                    <th>Коробок</th>
                </tr>
                </thead>
                <tbody>
                {print.register?.map(renderRow)}
                </tbody>
                <tfoot>
                <tr>
                    <th></th>
                    <th style={{textAlign: "right"}}>Итого:</th>
                    <th>{print.register?.reduce((l, r) => l + r.amount, 0)}</th>
                    <th>{print.register?.reduce((l, r) => l + r.boxes, 0)}</th>
                </tr>
                </tfoot>
            </table>
        </div>
    </div>;
}
