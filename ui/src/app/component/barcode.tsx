import JsBarcode, {Options} from "jsbarcode";
import React, {useEffect, useRef, useState} from "react";

interface Props {
    value: string;
    options?: Options;
}

const defaults: Options = {
    format: "EAN13",
    textMargin: 0,
    fontOptions: "bold",
}

export function Barcode({value, options}: Props) {
    const barcodeRef = useRef(null);
    const [invalid, setInvalid] = useState(false);

    useEffect(() => {
        const el = barcodeRef.current;
        const opts = {...defaults, options};
        if (el === null) {
            return;
        }
        try {
            JsBarcode(el, value, opts);
            setInvalid(false);
        } catch (e) {
            console.error(e);
            setInvalid(true);
        }
    }, [barcodeRef, value, options]);

    return <>
        {invalid ? <div>Некорректный штрих-код EAN-13: {value}</div> : null}
        <svg ref={barcodeRef}/>
    </>;
}
