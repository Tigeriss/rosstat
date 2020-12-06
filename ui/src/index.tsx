import React from 'react';
import ReactDOM from 'react-dom';
import {App} from './app/app';
import 'semantic-ui-css/semantic.min.css';
import "./index.css";
import {configure} from "mobx"

if (window.location.pathname.includes("/print/")) {
    document.body.className = "print";
}

document.addEventListener("wheel", function (event) {
    if ((document?.activeElement as any)?.type === "number") {
        (document.activeElement as any).blur();
    }
});

configure({
    enforceActions: "never",
})

ReactDOM.render(
    <App/>,
    document.getElementById('root')
);
