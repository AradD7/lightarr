import { useState } from "react";

export default function Commands(props) {
    const [commandsOpen, setCommandsOpen] = useState(false);

    const commands = ["Turn off", "Turn on", "Dim", "Change Color", "Change Temperature"];

    const getCommandValue = (command) => {
        const found = props.selectedCommands.find(c => c.command === command);
        return found?.value;
    };

    const isCommandSelected = (command) => {
        return props.selectedCommands.some(c => c.command === command);
    };

    const isCommandDisabled = (command) => {
        // Disable all commands except "Turn off" when "Turn off" is selected
        if (command !== "Turn off" && isCommandSelected("Turn off")) {
            return true;
        }
        // Disable "Turn off" when any other command is selected
        if (command === "Turn off" && props.selectedCommands.some(c => c.command !== "Turn off")) {
            return true;
        }
        // Disable "Change Color" when "Change Temperature" is selected
        if (command === "Change Color" && isCommandSelected("Change Temperature")) {
            return true;
        }
        // Disable "Change Temperature" when "Change Color" is selected
        if (command === "Change Temperature" && isCommandSelected("Change Color")) {
            return true;
        }
        return false;
    };

    const handleTempChange = (e) => {
        const value = parseInt(e.target.value);
        if (!isNaN(value)) {
            props.updateCommandValue("Change Temperature", value);
        }
    };

    const handleTempBlur = (e) => {
        const value = parseInt(e.target.value);
        if (isNaN(value) || value < 2200) {
            props.updateCommandValue("Change Temperature", 2200);
        } else if (value > 6500) {
            props.updateCommandValue("Change Temperature", 6500);
        }
    };

    return (
        <div
            className={props.isEmpty && props.selectedCommands.length === 0 ? "select-commands-empty" : "select-commands"}
            onBlur={(e) => {
                if (!e.currentTarget.contains(e.relatedTarget)) {
                    setCommandsOpen(false);
                }
            }}
            tabIndex={0}
        >
            <div
                className="select-commands-header"
                onClick={() => setCommandsOpen(prev => !prev)}
            >
                {props.selectedCommands.length === 0 ? "Select Commands..." : `(${props.selectedCommands.length}) Command${props.selectedCommands.length === 1 ? "" : "s"} Selected`} {commandsOpen ? '◀' : '▶'}
            </div>
            {commandsOpen && (
                <div className="select-commands-list">
                    {commands.map(command => (
                        <div key={command}>
                            <label
                                className="select-command-item"
                                style={{
                                    opacity: isCommandDisabled(command) ? 0.5 : 1,
                                    cursor: isCommandDisabled(command) ? 'not-allowed' : 'pointer'
                                }}
                            >
                                <input
                                    type="checkbox"
                                    checked={isCommandSelected(command)}
                                    onChange={() => props.addCommand(command)}
                                    disabled={isCommandDisabled(command)}
                                />
                                {command}
                            </label>

                            {/* Dim input */}
                            {command === "Dim" && isCommandSelected("Dim") && (
                                <div className="command-input-container">
                                    <label className="command-input-label">
                                        Brightness (1-100):
                                    </label>
                                    <input
                                        type="number"
                                        min="1"
                                        max="100"
                                        value={getCommandValue("Dim") || 50}
                                        onChange={(e) => {
                                            const val = Math.min(100, Math.max(1, parseInt(e.target.value) || 1));
                                            props.updateCommandValue("Dim", val);
                                        }}
                                        className="command-number-input"
                                    />
                                </div>
                            )}

                            {/* RGB inputs */}
                            {command === "Change Color" && isCommandSelected("Change Color") && (
                                <div className="command-input-container">
                                    <label className="command-input-label">RGB Values:</label>
                                    <div className="rgb-inputs">
                                        <div className="rgb-input-group">
                                            <label>R:</label>
                                            <input
                                                type="number"
                                                min="0"
                                                max="255"
                                                value={getCommandValue("Change Color")?.r || 255}
                                                onChange={(e) => {
                                                    const val = Math.min(255, Math.max(0, parseInt(e.target.value) || 0));
                                                    const currentRGB = getCommandValue("Change Color") || { r: 255, g: 255, b: 255 };
                                                    props.updateCommandValue("Change Color", { ...currentRGB, r: val });
                                                }}
                                                className="command-number-input rgb-input"
                                            />
                                        </div>
                                        <div className="rgb-input-group">
                                            <label>G:</label>
                                            <input
                                                type="number"
                                                min="0"
                                                max="255"
                                                value={getCommandValue("Change Color")?.g || 255}
                                                onChange={(e) => {
                                                    const val = Math.min(255, Math.max(0, parseInt(e.target.value) || 0));
                                                    const currentRGB = getCommandValue("Change Color") || { r: 255, g: 255, b: 255 };
                                                    props.updateCommandValue("Change Color", { ...currentRGB, g: val });
                                                }}
                                                className="command-number-input rgb-input"
                                            />
                                        </div>
                                        <div className="rgb-input-group">
                                            <label>B:</label>
                                            <input
                                                type="number"
                                                min="0"
                                                max="255"
                                                value={getCommandValue("Change Color")?.b || 255}
                                                onChange={(e) => {
                                                    const val = Math.min(255, Math.max(0, parseInt(e.target.value) || 0));
                                                    const currentRGB = getCommandValue("Change Color") || { r: 255, g: 255, b: 255 };
                                                    props.updateCommandValue("Change Color", { ...currentRGB, b: val });
                                                }}
                                                className="command-number-input rgb-input"
                                            />
                                        </div>
                                    </div>
                                </div>
                            )}

                            {/* Temperature input */}
                            {command === "Change Temperature" && isCommandSelected("Change Temperature") && (
                                <div className="command-input-container">
                                    <label className="command-input-label">
                                        Temperature (2200-6500K):
                                    </label>
                                    <input
                                        type="number"
                                        min="2200"
                                        max="6500"
                                        value={getCommandValue("Change Temperature") || 4000}
                                        onChange={handleTempChange}
                                        onBlur={handleTempBlur}
                                        className="command-number-input"
                                    />
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
