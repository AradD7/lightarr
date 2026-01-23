import { useState } from "react";

export default function Devices(props) {
    const [devicesOpen, setDevicesOpen] = useState(false);
    const [selectedDevices, setSelectedDevices] = useState([]);

    const toggleDevice = (deviceId) => {
        setSelectedDevices(prev =>
            prev.includes(deviceId)
                ? prev.filter(id => id !== deviceId)
                : [...prev, deviceId]
        );
        props.addDevice(deviceId)
    };

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
                {selectedDevices.length === 0 ? "Select Devices..." : `(${selectedDevices.length}) Device${selectedDevices.length === 1 ? "" : "s"} Selected`} {devicesOpen ? '◀' : '▶'}
            </div>
            {devicesOpen && (
                <div
                    className="select-devices-list"
                >
                    {
                        props.data ?
                            props.data.map(device => (
                                <label key={device.id} className="select-device-item">
                                    <input
                                        type="checkbox"
                                        checked={selectedDevices.includes(device.id)}
                                        onChange={() => toggleDevice(device.id)}
                                    />
                                    {device.name}
                                </label>
                            ))
                            : <section className="rules-empty-array">No Devices Saved</section>
                    }
                </div>
            )}

        </div>

    )
}
