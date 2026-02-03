import { useState } from "react";

export default function BulbsMac(props) {
    const [bulbsMacsOpen, setBulbsMacsOpen] = useState(false);

    return (
        <div
            className={props.isEmpty && props.selectedBulbsMacs.length === 0 ? "select-bulbsMacs-empty" : "select-bulbsMacs"}
            onBlur={(e) => {
                if (!e.currentTarget.contains(e.relatedTarget)) {
                    setBulbsMacsOpen(false);
                }
            }}
            tabIndex={0}
        >
            <div
                className="select-bulbsMacs-header"
                onClick={() => setBulbsMacsOpen(prev => !prev)}
            >
                {props.selectedBulbsMacs.length === 0 ? "Select Bulbs..." : `(${props.selectedBulbsMacs.length}) Bulb${props.selectedBulbsMacs.length === 1 ? "" : "s"} Selected`} {bulbsMacsOpen ? '◀' : '▶'}
            </div>
            {bulbsMacsOpen && (
                <div
                    className="select-bulbsMacs-list"
                >
                    {props.data && props.data.length > 0 ?
                        <>
                            <label className="select-bulbsMacs-item">
                                <input
                                    type="checkbox"
                                    checked={props.selectedBulbsMacs.length === props.data.length}
                                    onChange={() => props.addBulbsMac("toggleAll")}
                                />
                                Select All
                            </label>
                            {props.data.map(bulbsMac => (
                                <label key={bulbsMac.mac} className="select-bulbsMacs-item">
                                    <input
                                        type="checkbox"
                                        checked={props.selectedBulbsMacs.includes(bulbsMac.mac)}
                                        onChange={() => props.addBulbsMac(bulbsMac.mac)}
                                    />
                                    {bulbsMac.name}
                                </label>
                            ))}
                        </>
                        : <section className="rules-empty-array">No Bulbs Found. Head to Bulbs tab.</section>
                    }
                </div>
            )}
        </div>
    )
}
