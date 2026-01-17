import { useState } from "react";

export default function Commands(props) {
    const [commandsOpen, setCommandsOpen] = useState(false);
    const [selectedCommands, setSelectedCommands] = useState([]);

    const commands = ["Turn off", "Turn on", "Dim", "Change Color", "Change Temperature"];

    const getCommandValue = (command) => {
        const found = selectedCommands.find(c => c.command === command);
        return found?.value;
    };

    const isCommandSelected = (command) => {
        return selectedCommands.some(c => c.command === command);
    };

    const toggleCommand = (command) => {
        setSelectedCommands(prev => {
            const isSelected = prev.some(c => c.command === command);

            if (isSelected) {
                return prev.filter(c => c.command !== command);
            }

            // If selecting "Turn off", clear all other selections
            if (command === "Turn off") {
                return [{ command: "Turn off" }];
            }

            // If selecting anything else, remove "Turn off"
            let filtered = prev.filter(c => c.command !== "Turn off");

            // Handle Change Color and Change Temperature mutual exclusivity
            if (command === "Change Color") {
                filtered = filtered.filter(c => c.command !== "Change Temperature");
            }
            if (command === "Change Temperature") {
                filtered = filtered.filter(c => c.command !== "Change Color");
            }

            // Add new command with default values
            let newCommand;
            if (command === "Dim") {
                newCommand = { command, value: 50 };
            } else if (command === "Change Color") {
                newCommand = { command, value: { r: 255, g: 255, b: 255 } };
            } else if (command === "Change Temperature") {
                newCommand = { command, value: 4000 };
            } else {
                newCommand = { command };
            }

            return [...filtered, newCommand];
        });
        props.addCommand(command);
    };

    const updateCommandValue = (command, value) => {
        setSelectedCommands(prev =>
            prev.map(c => c.command === command ? { ...c, value } : c)
        );
    };

    const isCommandDisabled = (command) => {
        // Disable all commands except "Turn off" when "Turn off" is selected
        if (command !== "Turn off" && isCommandSelected("Turn off")) {
            return true;
        }
        // Disable "Turn off" when any other command is selected
        if (command === "Turn off" && selectedCommands.some(c => c.command !== "Turn off")) {
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
            updateCommandValue("Change Temperature", value);
        }
    };

    const handleTempBlur = (e) => {
        const value = parseInt(e.target.value);
        if (isNaN(value) || value < 2200) {
            updateCommandValue("Change Temperature", 2200);
        } else if (value > 6500) {
            updateCommandValue("Change Temperature", 6500);
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
                                    checked={isCommandSelected(command)}
                                    onChange={() => toggleCommand(command)}
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
                                            updateCommandValue("Dim", val);
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
                                                    updateCommandValue("Change Color", { ...currentRGB, r: val });
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
                                                    updateCommandValue("Change Color", { ...currentRGB, g: val });
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
                                                    updateCommandValue("Change Color", { ...currentRGB, b: val });
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
