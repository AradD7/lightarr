import { useState } from "react";

export default function BulbsMac(props) {
    const [bulbsMacsOpen, setBulbsMacsOpen] = useState(false);
    const [selectedBulbsMacs, setSelectedBulbsMacs] = useState([]);

    const toggleBulbsMac = (bulbsMac) => {
        setSelectedBulbsMacs(prev =>
            prev.includes(bulbsMac)
                ? prev.filter(mac => mac !== bulbsMac)
                : [...prev, bulbsMac]
        );
        props.addBulbsMac(bulbsMac)
    };

    return (
        <div
            className="select-bulbsMacs"
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
                {selectedBulbsMacs.length === 0 ? "Select Bulbs..." : `(${selectedBulbsMacs.length}) Bulbs${selectedBulbsMacs.length === 1 ? "" : "s"} Selected`} {bulbsMacsOpen ? '◀' : '▶'}
            </div>
            {bulbsMacsOpen && (
                <div
                    className="select-bulbsMacs-list"
                >
                    {props.data.map(bulbsMac => (
                        <label key={bulbsMac.mac} className="select-bulbsMacs-item">
                            <input
                                type="checkbox"
                                checked={selectedBulbsMacs.includes(bulbsMac.mac)}
                                onChange={() => toggleBulbsMac(bulbsMac.mac)}
                            />
                            {bulbsMac.name}
                        </label>
                    ))}
                </div>
            )}

        </div>

    )
}
