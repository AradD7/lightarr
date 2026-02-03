import { useState } from "react";

export default function Devices(props) {
    const [devicesOpen, setDevicesOpen] = useState(false);

    return (
        <div
            className="select-devices"
            onBlur={(e) => {
                if (!e.currentTarget.contains(e.relatedTarget)) {
                    setDevicesOpen(false);
                }
            }}
            tabIndex={0}
        >
            <div
                className="select-devices-header"
                onClick={() => setDevicesOpen(prev => !prev)}
            >
                {props.selectedDevices.length === 0 ? "Select Devices..." : `(${props.selectedDevices.length}) Device${props.selectedDevices.length === 1 ? "" : "s"} Selected`} {devicesOpen ? '◀' : '▶'}
            </div>
            {devicesOpen && (
                <div
                    className="select-devices-list"
                >
                    {props.data && props.data.length > 0 ?
                        <>
                            <label className="select-device-item">
                                <input
                                    type="checkbox"
                                    checked={props.selectedDevices.length === props.data.length}
                                    onChange={() => props.addDevice("toggleAll")}
                                />
                                Select All
                            </label>
                            {props.data.map(device => (
                                <label key={device.id} className="select-device-item">
                                    <input
                                        type="checkbox"
                                        checked={props.selectedDevices.includes(device.id)}
                                        onChange={() => props.addDevice(device.id)}
                                    />
                                    {device.name}
                                </label>
                            ))}
                        </>
                        : <section className="rules-empty-array">No Devices Saved</section>
                    }
                </div>
            )}
        </div>
    )
}
