import { useState } from "react";

export default function Commands() {
    const [commandsOpen, setCommandsOpen] = useState(false);
    const [selectedCommands, setSelectedCommands] = useState([]);
    const [dimValue, setDimValue] = useState(50);
    const [rgbValues, setRgbValues] = useState({ r: 255, g: 255, b: 255 });
    const [tempValue, setTempValue] = useState(4000);

    const commands = ["Turn off", "Turn on", "Dim", "Change Color", "Change Temperature"];

    const toggleCommand = (command) => {
        setSelectedCommands(prev => {
            if (prev.includes(command)) {
                return prev.filter(cmd => cmd !== command);
            }

            // If selecting "Turn off", clear all other selections
            if (command === "Turn off") {
                return ["Turn off"];
            }

            // If selecting anything else, remove "Turn off"
            const filtered = prev.filter(cmd => cmd !== "Turn off");

            // Handle Change Color and Change Temperature mutual exclusivity
            if (command === "Change Color") {
                return [...filtered.filter(cmd => cmd !== "Change Temperature"), command];
            }
            if (command === "Change Temperature") {
                return [...filtered.filter(cmd => cmd !== "Change Color"), command];
            }

            return [...filtered, command];
        });
    };

    const isCommandDisabled = (command) => {
        // Disable all commands except "Turn off" when "Turn off" is selected
        if (command !== "Turn off" && selectedCommands.includes("Turn off")) {
            return true;
        }
        // Disable "Turn off" when any other command is selected
        if (command === "Turn off" && selectedCommands.some(cmd => cmd !== "Turn off")) {
            return true;
        }
        // Disable "Change Color" when "Change Temperature" is selected
        if (command === "Change Color" && selectedCommands.includes("Change Temperature")) {
            return true;
        }
        // Disable "Change Temperature" when "Change Color" is selected
        if (command === "Change Temperature" && selectedCommands.includes("Change Color")) {
            return true;
        }
        return false;
    };

    const handleTempChange = (e) => {
        const value = parseInt(e.target.value);
        if (!isNaN(value)) {
            setTempValue(value);
        }
    };

    const handleTempBlur = (e) => {
        const value = parseInt(e.target.value);
        if (isNaN(value) || value < 2200) {
            setTempValue(2200);
        } else if (value > 6500) {
            setTempValue(6500);
        }
    };

    return (
        <div
            className="select-commands"
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
                Select commands... {commandsOpen ? '◀' : '▶'}
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
                                    checked={selectedCommands.includes(command)}
                                    onChange={() => toggleCommand(command)}
                                    disabled={isCommandDisabled(command)}
                                />
                                {command}
                            </label>

                            {/* Dim input */}
                            {command === "Dim" && selectedCommands.includes("Dim") && (
                                <div className="command-input-container">
                                    <label className="command-input-label">
                                        Brightness (1-100):
                                    </label>
                                    <input
                                        type="number"
                                        min="1"
                                        max="100"
                                        value={dimValue}
                                        onChange={(e) => setDimValue(Math.min(100, Math.max(1, parseInt(e.target.value) || 1)))}
                                        className="command-number-input"
                                    />
                                </div>
                            )}

                            {/* RGB inputs */}
                            {command === "Change Color" && selectedCommands.includes("Change Color") && (
                                <div className="command-input-container">
                                    <label className="command-input-label">RGB Values:</label>
                                    <div className="rgb-inputs">
                                        <div className="rgb-input-group">
                                            <label>R:</label>
                                            <input
                                                type="number"
                                                min="0"
                                                max="255"
                                                value={rgbValues.r}
                                                onChange={(e) => setRgbValues(prev => ({
                                                    ...prev,
                                                    r: Math.min(255, Math.max(0, parseInt(e.target.value) || 0))
                                                }))}
                                                className="command-number-input rgb-input"
                                            />
                                        </div>
                                        <div className="rgb-input-group">
                                            <label>G:</label>
                                            <input
                                                type="number"
                                                min="0"
                                                max="255"
                                                value={rgbValues.g}
                                                onChange={(e) => setRgbValues(prev => ({
                                                    ...prev,
                                                    g: Math.min(255, Math.max(0, parseInt(e.target.value) || 0))
                                                }))}
                                                className="command-number-input rgb-input"
                                            />
                                        </div>
                                        <div className="rgb-input-group">
                                            <label>B:</label>
                                            <input
                                                type="number"
                                                min="0"
                                                max="255"
                                                value={rgbValues.b}
                                                onChange={(e) => setRgbValues(prev => ({
                                                    ...prev,
                                                    b: Math.min(255, Math.max(0, parseInt(e.target.value) || 0))
                                                }))}
                                                className="command-number-input rgb-input"
                                            />
                                        </div>
                                    </div>
                                </div>
                            )}

                            {/* Temperature input */}
                            {command === "Change Temperature" && selectedCommands.includes("Change Temperature") && (
                                <div className="command-input-container">
                                    <label className="command-input-label">
                                        Temperature (2200-6500K):
                                    </label>
                                    <input
                                        type="number"
                                        min="2200"
                                        max="6500"
                                        value={tempValue}
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
